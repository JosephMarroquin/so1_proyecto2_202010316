# Usa una imagen de Go como base
FROM golang:latest

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia el código fuente al contenedor
COPY . .

# Descarga e instala las dependencias
RUN go mod download

# Compila la aplicación
RUN go build -o main .

# Expone el puerto 8080 (si tu aplicación Go es un servidor HTTP)
EXPOSE 3000

# Ejecuta la aplicación al iniciar el contenedor
CMD ["./main"]
