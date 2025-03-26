# Etapa 1: Build da aplicação
FROM golang:1.24 AS builder

WORKDIR /app

# Copia os arquivos da aplicação para o contêiner
COPY . .

# Compila a aplicação (binário leve para Linux)
RUN CGO_ENABLED=0 GOOS=linux go build -o stress-test ./cmd/stresstest

# Etapa 2: Imagem final mínima
FROM scratch

# Copia apenas o binário para a imagem final
COPY --from=builder /app/stress-test /stress-test

# Define o binário como ponto de entrada
ENTRYPOINT ["/stress-test"]