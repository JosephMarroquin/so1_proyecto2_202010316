package main

import (
	"context"
	"fmt"
	pb "grpcServer/server"
	"log"
	"net"

	"os"

	"google.golang.org/grpc"

	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type server struct {
	pb.UnimplementedGetInfoServer
}

const (
	port = ":3000"
)

type Data struct {
	Album  string
	Year   string
	Artist string
	Ranked string
}

func sendMessage(topic, message string) error {
	// Load Kafka brokers from environment variable
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		return fmt.Errorf("KAFKA_BROKERS environment variable is not set")
	}

	// Create a new Kafka writer
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBrokers},
		Topic:   topic,
	})

	defer w.Close()

	// Produce message to Kafka topic
	err := w.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})

	if err != nil {
		return fmt.Errorf("Error sending message: %s", err)
	}

	return nil
}

func (s *server) ReturnInfo(ctx context.Context, in *pb.RequestId) (*pb.ReplyInfo, error) {
	//fmt.Println("Recibí de cliente: ", in.GetArtist())
	data := Data{
		Year:   in.GetYear(),
		Album:  in.GetAlbum(),
		Artist: in.GetArtist(),
		Ranked: in.GetRanked(),
	}
	messageBytes, err2 := json.Marshal(data)
	if err2 != nil {
		fmt.Println("Error al convertir a JSON:", err2)
		//return
	}
	// Convertir la cadena JSON a una cadena de texto
	message := string(messageBytes)

	topic := "topic-sopes1"
	err := sendMessage(topic, message)
	if err != nil {
		fmt.Printf("Error al enviar el mensaje: %s\n", err)
		//return
	}

	fmt.Println("Mensaje enviado exitosamente")
	//fmt.Println(data)
	// insertRedis(data)
	return &pb.ReplyInfo{Info: "Hola cliente, recibí el album"}, nil
}

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	pb.RegisterGetInfoServer(s, &server{})

	// redisConnect()

	if err := s.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}
