# --- ETAPA 1: Build (Usamos la última versión de Go) ---
# Usamos una etiqueta específica de Alpine reciente para evitar sorpresas
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instalamos git y certificados actualizados
RUN apk update && apk add --no-cache git ca-certificates tzdata

# Gestión de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el código
COPY . .

# Compilación ESTÁTICA (Vital para Distroless)
# -ldflags="-s -w": Quita información de depuración para bajar el peso
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# --- ETAPA 2: Producción (Google Distroless) ---
# Usamos 'static-debian12' que es ultra ligera y segura
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copiamos el binario de la etapa anterior
COPY --from=builder /app/main .

# Copiamos la zona horaria (opcional, pero útil para logs correctos)
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Si tu app usa un .env, descomenta esto. 
# PERO RECOMENDACIÓN: Inyecta las variables por docker-compose, no copies el archivo.
# COPY .env .

# Puerto (cámbialo si tu main.go usa otro, ej: 3000)
EXPOSE 8080

# Usuario no-root (seguridad extra que distroless ya maneja, pero definimos user id)
USER nonroot:nonroot

CMD ["/main"]