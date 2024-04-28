package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"

	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"net/url"
	"time"
)

type DatosVotacion struct {
	Album  string `json:"Album"`
	Year   string `json:"Year"`
	Artist string `json:"Artist"`
	Ranked string `json:"Ranked"`
}

func redis_conn(albumKey string) {
	// Crea un cliente Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",            // Utiliza el nombre del servicio y el puerto
		Password: os.Getenv("REDIS_AUTH"), // si tienes autenticación, especifica la contraseña aquí
		DB:       0,                       // Número de la base de datos que deseas seleccionar
	})

	// Crea un contexto
	ctx := context.Background()

	// Verificar si ya existe en Redis
	exists, errExists := rdb.HExists(ctx, "datos", albumKey).Result()
	if errExists != nil {
		fmt.Println("Error al verificar la existencia del álbum en Redis:", errExists)
		return
	}

	// Si el álbum no existe, agregarlo a Redis con un contador de votos inicial de 0
	if !exists {
		err := rdb.HSet(ctx, "datos", albumKey, 0).Err()
		if err != nil {
			fmt.Println("Error al agregar el álbum a Redis:", err)
			return
		}
	}

	// Incrementar el contador de votos para el álbum
	_, errExists = rdb.HIncrBy(ctx, "datos", albumKey, 1).Result()
	if errExists != nil {
		fmt.Println("Error al incrementar el contador de votos:", errExists)
		return
	}

	// Obtener el contador de votos actualizado para el álbum
	votos, err := rdb.HGet(ctx, "datos", albumKey).Result()
	if err != nil {
		fmt.Println("Error al obtener el contador de votos:", err)
		return
	}

	fmt.Printf("El contador de votos para %s es %s\n", albumKey, votos)

	// Espera un poco antes de cerrar la conexión
	time.Sleep(1 * time.Second)

	// Cierra la conexión al final
	err = rdb.Close()
	if err != nil {
		panic(err)
	}
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Message   string    `json:"message,omitempty"`
}

func mongo_conn(message string) {
	// Establecer la conexión con la base de datos MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://" + os.Getenv("MONGO_INITDB_ROOT_USERNAME") + ":" + url.QueryEscape(os.Getenv("MONGO_INITDB_ROOT_PASSWORD")) + "@mongodb-service:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Seleccionar la base de datos y la colección
	collection := client.Database("proyecto2").Collection("logs")

	entry := LogEntry{
		Timestamp: time.Now(),
		Message:   message,
	}

	// Insertar un documento JSON
	_, err = collection.InsertOne(context.Background(), entry)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Documento insertado exitosamente.")
}

func main() {
	// Read environment variables
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	// Initialize Kafka reader
	config := kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   kafkaTopic,
		GroupID: "test-group",
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	fmt.Printf("Listening to topic: %s\n", kafkaTopic)

	// Handle OS signals to gracefully shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Start reading messages
	for {
		select {
		case <-sigchan:
			fmt.Println("Shutting down...")
			return
		default:
			msg, err := reader.FetchMessage(context.Background())
			if err != nil {
				fmt.Printf("Error fetching message: %v\n", err)
				continue
			}
			fmt.Printf("Message from topic %s: %s\n", msg.Topic, msg.Value)

			//convertir string a json del msg.value
			var datosVotacion DatosVotacion

			errV := json.Unmarshal([]byte(msg.Value), &datosVotacion)
			if errV != nil {
				fmt.Println("Error", errV)
				return
			}

			//fmt.Println("Nombre: ", datosVotacion.Artist)
			//fmt.Println("album: ", datosVotacion.Album)
			//fmt.Println("year: ", datosVotacion.Year)
			//fmt.Println("rank: ", datosVotacion.Ranked)

			albumKey := strings.Join([]string{datosVotacion.Artist, datosVotacion.Album, datosVotacion.Year}, "_")
			//fmt.Println("key: ", albumKey)

			//contador para redis
			go redis_conn(albumKey)
			//logs para db de mongo
			go mongo_conn(albumKey)
			err = reader.CommitMessages(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error committing message: %v\n", err)
			}
		}
	}
}
