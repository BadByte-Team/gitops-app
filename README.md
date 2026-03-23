# gitops-app — Código Fuente del Curso GitOps

Aplicación web Go con autenticación JWT, panel de administración y frontend HTML estático. Diseñada como proyecto práctico del curso de GitOps con arquitectura híbrida gratuita (Jenkins local + K3s en AWS Free Tier).

---

## Stack tecnológico

| Capa | Tecnología |
|---|---|
| Backend | Go 1.22 — `net/http` estándar |
| Base de datos | MySQL 8.0 |
| Auth | JWT (golang-jwt/jwt/v5) |
| Frontend | HTML + CSS + JS vanilla |
| CI | Jenkins (Docker local) |
| CD | ArgoCD en K3s |

---

## Estructura

```
gitops-app/
├── cmd/api/main.go           # Punto de entrada
├── internal/
│   ├── auth/jwt.go           # Generación y validación JWT + middlewares
│   ├── handlers/handlers.go  # Handlers HTTP
│   ├── models/models.go      # Structs del dominio
│   └── repository/
│       ├── db.go             # Conexión a MySQL
│       └── user.go           # Repositorios de usuarios y módulos
├── frontend/
│   ├── index.html            # Login
│   ├── dashboard.html        # Panel del estudiante
│   └── admin.html            # Panel de administración
├── init.sql                  # Schema para pruebas locales
├── Dockerfile                # Multi-stage build (Go 1.22 + Alpine 3.19)
├── docker-compose.yml        # App + MySQL para desarrollo local
└── Jenkinsfile               # Pipeline CI de 6 stages
```

---

## Desarrollo local

```bash
# 1. Clonar el repositorio
git clone https://github.com/TU_USUARIO/gitops-app.git
cd gitops-app

# 2. Levantar app + MySQL
docker compose up -d

# 3. Verificar
curl http://localhost:8080/api/health
# {"success":true,"message":"ok"}

# 4. Login de prueba
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

---

## Variables de entorno

| Variable | Default | Descripción |
|---|---|---|
| `DB_HOST` | `localhost` | Host de MySQL |
| `DB_PORT` | `3306` | Puerto de MySQL |
| `DB_USER` | `curso_app` | Usuario de la BD |
| `DB_PASSWORD` | `C4rs0_S3cur3_P@ss!` | Contraseña |
| `DB_NAME` | `curso_db` | Nombre de la BD |
| `JWT_SECRET` | `dev_secret_...` | Secreto para firmar JWT |
| `PORT` | `8080` | Puerto del servidor HTTP |

---

## Pipeline CI (Jenkinsfile)

El pipeline ejecuta 6 stages en orden:

1. **Checkout** — clona el repo y genera `BUILD_TAG = N-GIT_HASH`
2. **SonarQube Analysis** — analiza calidad del código Go
3. **Quality Gate** — espera el resultado de SonarQube (falla si no pasa)
4. **Docker Build** — construye la imagen con el `BUILD_TAG`
5. **Trivy Scan** — escanea vulnerabilidades HIGH/CRITICAL
6. **Docker Push** — sube la imagen a Docker Hub
7. **Deploy to GitOps Repo** — actualiza `deployment.yaml` en `gitops-infra` con `sed`
8. **Cleanup** — limpia imágenes locales

---

## Credenciales en Jenkins (EP34)

| ID | Tipo | Uso |
|---|---|---|
| `dockerhub-id` | Username/password | Push a Docker Hub |
| `github-token-id` | Secret text | Clonar y hacer push a `gitops-infra` |
| `sonarqube-token` | Secret text | Autenticar con SonarQube |
