# Etapa 1: Build
FROM golang:1.25-alpine AS builder
# Instala git e certificados necessários
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copia apenas os arquivos de módulo primeiro
COPY go.mod go.sum ./

# Tenta baixar as dependências com um "retry" manual ou ignorando somas se necessário
# mas o tidy feito antes deve resolver
RUN go mod download

# Agora copia o resto do código
COPY . .

# Compila os binários (CGO_ENABLED=0 garante portabilidade no Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o seeder cmd/seeder/*.go

FROM alpine:latest

# INSTALE O TZDATA AQUI
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/seeder .
COPY --from=builder /app/config ./config
COPY --from=builder /app/dados_importacao ./dados_importacao

EXPOSE 8080
CMD ["./main"]
