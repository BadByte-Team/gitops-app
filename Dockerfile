# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git ca-certificates tzdata

# Directorio de trabajo
WORKDIR /build

# Descargar dependencias primero (capa cacheable)
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar el binario estático
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o curso-gitops \
    ./cmd/api/main.go

# ── Stage 2: Runtime ───────────────────────────────────────────────────────────
FROM alpine:3.19

# Seguridad: no correr como root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Certificados para HTTPS y timezone
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copiar solo el binario y el frontend del stage de build
COPY --from=builder /build/curso-gitops .
COPY --from=builder /build/frontend ./frontend

# Cambiar propietario
RUN chown -R appuser:appgroup /app

USER appuser

# Puerto de la app
EXPOSE 8080

# Health check para Kubernetes readiness probe
HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

CMD ["./curso-gitops"]
