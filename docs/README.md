# gitops-app — Código Fuente del Curso GitOps

Repositorio de la **aplicación Go** del curso. Contiene el código fuente, el Dockerfile, el docker-compose para pruebas locales y el Jenkinsfile del pipeline CI.

> **Regla GitOps:** este repositorio es tocado por los **desarrolladores**. Jenkins lo clona para construir y analizar. ArgoCD **nunca** lee este repositorio — ArgoCD solo observa `gitops-infra`.

---

## ¿Qué hace esta aplicación?

Una plataforma web de curso online con:

- **Login con JWT** — el usuario recibe un token firmado con su rol (`admin` o `student`)
- **Dashboard del estudiante** — lista los módulos del curso visibles
- **Panel de administración** — lista todos los usuarios y todos los módulos (incluyendo los ocultos)
- **API REST** con rutas protegidas por middleware JWT
- **Frontend HTML estático** servido por la misma app Go en el puerto 8080

---

## Stack técnico

| Capa | Tecnología | Versión |
|---|---|---|
| Lenguaje | Go | 1.22 |
| Base de datos | MySQL | 8.0 |
| Auth | JWT | golang-jwt/jwt v5 |
| Hashing contraseñas | bcrypt | golang.org/x/crypto |
| Frontend | HTML + CSS + JS vanilla | sin frameworks |
| Imagen base (build) | golang:1.22-alpine | — |
| Imagen base (runtime) | alpine | 3.19 |

---

## Estructura completa del repositorio

```
gitops-app/
│
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada
│                                # Registra las rutas HTTP, conecta a MySQL
│                                # con reintentos (10 intentos × 3s), arranca en :8080
│
├── internal/
│   ├── auth/
│   │   └── jwt.go               # GenerateToken()  — crea JWT firmado con HS256, expira en 24h
│   │                            # ValidateToken()  — valida firma y expiración
│   │                            # RequireAuth()    — middleware: extrae y valida el Bearer token
│   │                            # RequireAdmin()   — middleware: verifica que el rol sea "admin"
│   │
│   ├── handlers/
│   │   └── handlers.go          # Login()       — POST /api/login
│   │                            # GetModules()  — GET  /api/modules (autenticado)
│   │                            # GetUsers()    — GET  /api/users   (solo admin)
│   │                            # Health()      — GET  /api/health  (público, para K8s probes)
│   │
│   ├── models/
│   │   └── models.go            # User, Module, Episode (structs del dominio)
│   │                            # LoginRequest, LoginResponse, APIResponse (structs de la API)
│   │
│   └── repository/
│       ├── db.go                # NewDB() — abre la conexión MySQL usando variables de entorno
│       │                        # Configura pool: 10 conexiones máximas, 5 idle
│       └── user.go              # UserRepository  — FindByUsername(), Create(), ListAll()
│                                # ModuleRepository — ListVisible(), ListAll()
│
├── frontend/
│   ├── index.html               # Página de login
│   │                            # Llama POST /api/login, guarda el JWT en localStorage
│   │                            # Redirige a /admin.html si es admin, /dashboard.html si es student
│   │
│   ├── dashboard.html           # Panel del estudiante
│   │                            # Llama GET /api/modules con el Bearer token
│   │                            # Muestra la lista de módulos visibles
│   │
│   └── admin.html               # Panel de administración
│                                # Llama GET /api/users y GET /api/modules
│                                # Muestra tablas con usuarios y estado de todos los módulos
│
├── init.sql                     # Schema completo de la BD + datos de prueba
│                                # Crea tablas: users, modules, episodes
│                                # Inserta usuario admin (contraseña: admin123)
│                                # Inserta 11 módulos del curso
│                                # Usado por docker-compose.yml localmente
│                                # En Kubernetes lo gestiona mysql-configmap.yaml
│
├── go.mod                       # Módulo Go — declara las dependencias del proyecto
│
├── Dockerfile                   # Multi-stage build
│                                # Stage 1 "builder": compila en golang:1.22-alpine
│                                # Stage 2 "runtime": copia solo el binario a alpine:3.19
│                                # Resultado: imagen ~20MB (vs ~300MB con imagen estándar de Go)
│                                # El contenedor corre como usuario no-root (appuser)
│
├── docker-compose.yml           # Levanta app + MySQL para desarrollo local
│                                # La app espera a que MySQL esté healthy antes de arrancar
│                                # Incluye health check de MySQL con mysqladmin ping
│
├── Jenkinsfile                  # Pipeline CI con 8 stages
│                                # Ver sección "Stages del pipeline" más abajo
│
├── .gitignore                   # Excluye binarios Go, .env, *.pem, logs
└── README.md                    # Este archivo
```

---

## En qué episodios se usa cada archivo

| Archivo / Directorio | Episodio | Qué se hace exactamente |
|---|---|---|
| `internal/` + `frontend/` + `init.sql` | **EP09 — Introducción a Docker** | Se presenta la app que se va a containerizar. El alumno la ve correr con `go run` antes de dockerizarla |
| `Dockerfile` | **EP10 — Dockerfile multi-stage** | Se escribe el Dockerfile explicando por qué dos stages, qué es `CGO_ENABLED=0`, por qué Alpine y no Ubuntu |
| `docker-compose.yml` | **EP11 — Docker Compose** | Se levanta la app + MySQL por primera vez. Se explica `depends_on`, redes y puertos |
| `docker-compose.yml` + `init.sql` | **EP12 — Docker Compose avanzado** | Se añaden health checks, volúmenes persistentes para MySQL, variables de entorno |
| `.gitignore` | **EP07 — Gitflow y estructura de repos** | Se crea el `.gitignore` para `gitops-app`. Se explica por qué los binarios y `.pem` no van al repo |
| `Jenkinsfile` (estructura básica) | **EP35 — Primer Jenkinsfile** | Se escribe el esqueleto: `agent`, `tools`, `environment`, `stages`, `post`. Sin stages reales todavía |
| `Jenkinsfile` (stages 1–4 + Deploy GitOps) | **EP36 — Pipeline CI con GitOps** | Se completa el Jenkinsfile con Checkout, Docker Build, Docker Push y el stage `Deploy to GitOps Repo` con `sed` |
| `Jenkinsfile` + montaje de Trivy | **EP42 — Trivy local** | Se añade el stage `Trivy Scan` al Jenkinsfile. Se verifica que Jenkins puede ejecutar `trivy` |
| `Jenkinsfile` + SonarQube | **EP44 — SonarQube + Jenkins** | Se añade el stage `SonarQube Analysis` y `Quality Gate`. Se configura el webhook de SonarQube → Jenkins |
| `Jenkinsfile` completo | **EP45 — Stages de seguridad** | Se revisa el orden de los 8 stages y se explica el principio "fail fast and early" |
| Todo el repositorio | **EP48 — Pipeline en acción** | Se hace un push real, se observa el pipeline correr, se verifica la imagen en Docker Hub y el commit en `gitops-infra` |

---

## Endpoints de la API

| Método | Ruta | Autenticación | Descripción |
|---|---|---|---|
| `GET` | `/api/health` | No | Readiness check para los probes de Kubernetes. Devuelve `{"success":true,"message":"ok"}` |
| `POST` | `/api/login` | No | Body: `{"username":"...","password":"..."}`. Devuelve `{"token":"eyJ...","role":"admin"}` |
| `GET` | `/api/modules` | JWT Bearer | Estudiante: solo módulos con `is_hidden=false`. Admin: todos los módulos |
| `GET` | `/api/users` | JWT Bearer (rol admin) | Lista todos los usuarios con id, username, role y fecha de registro |

---

## Variables de entorno

| Variable | Valor en desarrollo local | Valor en Kubernetes | Descripción |
|---|---|---|---|
| `DB_HOST` | `mysql-db` (docker-compose) | `mysql-svc` (DNS del Service K8s) | Host de MySQL |
| `DB_PORT` | `3306` | `3306` | Puerto de MySQL |
| `DB_USER` | `curso_app` | Desde Secret `db-credentials` | Usuario de la base de datos |
| `DB_PASSWORD` | `C4rs0_S3cur3_P@ss!` | Desde Secret `db-credentials` | Contraseña de la base de datos |
| `DB_NAME` | `curso_db` | `curso_db` | Nombre de la base de datos |
| `JWT_SECRET` | `dev_secret_change_in_production` | Desde Secret `app-secrets` | Clave para firmar los tokens JWT |
| `PORT` | `8080` | `8080` | Puerto en el que escucha el servidor HTTP |

---

## Credenciales que debes configurar en Jenkins (EP34)

Estos IDs deben coincidir **exactamente** con los que aparecen en el `Jenkinsfile`. Si los creas con nombres diferentes, el pipeline falla.

| ID en Jenkins | Tipo | Dónde obtenerlo | Para qué sirve |
|---|---|---|---|
| `dockerhub-id` | Username with password | Docker Hub → Account Settings → Security | Hacer `docker login` y `docker push` |
| `github-token-id` | Secret text (PAT) | GitHub → Settings → Developer settings → Personal access tokens → Classic → scope: `repo` | Clonar `gitops-infra` y hacer `git push` desde el stage `Deploy to GitOps Repo` |
| `sonarqube-token` | Secret text | SonarQube → My Account → Security → Generate token | Autenticar el análisis de código con SonarQube |

> ⚠️ El PAT de GitHub **debe tener el scope `repo`** (no solo `read`). Sin ese permiso, el `git push` a `gitops-infra` falla con error 403.

---

## Stages del pipeline (Jenkinsfile)

```
1. Checkout
   └─ Clona gitops-app
   └─ Genera BUILD_TAG = "N-GITHASH" (ej: "5-a3b8d1c")

2. SonarQube Analysis
   └─ Ejecuta sonar-scanner sobre el código Go
   └─ Excluye vendor/, node_modules/, frontend/

3. Quality Gate
   └─ Espera hasta 5 minutos el resultado de SonarQube
   └─ Si no pasa → el pipeline aborta aquí (no construye imagen)

4. Docker Build
   └─ docker build -t TU_USUARIO/curso-gitops:5-a3b8d1c .
   └─ También tagea como :latest

5. Trivy Scan
   └─ Escanea la imagen recién construida
   └─ Reporta vulnerabilidades HIGH y CRITICAL
   └─ --exit-code 0 → reporta pero no falla el pipeline

6. Docker Push
   └─ docker login con dockerhub-id
   └─ Sube el tag versionado (5-a3b8d1c) y latest a Docker Hub

7. Deploy to GitOps Repo ← el corazón del patrón GitOps
   └─ git clone gitops-infra usando github-token-id
   └─ sed actualiza deployment.yaml:
      ANTES:  image: TU_USUARIO/curso-gitops:4-xyz
      DESPUÉS: image: TU_USUARIO/curso-gitops:5-a3b8d1c
   └─ git commit -m "ci: deploy version 5-a3b8d1c from Jenkins"
   └─ git push → ArgoCD detecta este commit

8. Cleanup
   └─ Elimina imágenes locales para liberar disco
```

---

## Qué cambiar antes de usar este repositorio

Busca y reemplaza las siguientes dos cadenas en todos los archivos:

| Placeholder | Reemplazar con | Archivos afectados |
|---|---|---|
| `TU_USUARIO_DOCKERHUB` | Tu usuario de Docker Hub | `Jenkinsfile` |
| `TU_USUARIO_GITHUB` | Tu usuario de GitHub | `Jenkinsfile` |
| `TU_USUARIO` | Tu usuario (ambos) | `go.mod` + todos los `import` de `internal/` |

Ejemplo con `sed`:

```bash
# Reemplazar en Jenkinsfile
sed -i 's/TU_USUARIO_DOCKERHUB/johndoe/g' Jenkinsfile
sed -i 's/TU_USUARIO_GITHUB/johndoe/g' Jenkinsfile

# Reemplazar en go.mod y todos los archivos Go
find . -name "*.go" -o -name "go.mod" | xargs sed -i 's/TU_USUARIO/johndoe/g'
```

---

## Cómo probar en local (sin Kubernetes)

```bash
# 1. Clonar
git clone https://github.com/TU_USUARIO/gitops-app.git
cd gitops-app

# 2. Levantar app + MySQL
docker compose up -d

# 3. Verificar que los contenedores están corriendo
docker compose ps
# NAME                   STATUS
# curso-gitops-app       Up
# curso-gitops-mysql     Up (healthy)

# 4. Health check
curl http://localhost:8080/api/health
# {"success":true,"message":"ok"}

# 5. Login con el usuario admin de prueba
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
# {"success":true,"data":{"token":"eyJ...","role":"admin"}}

# 6. Guardar el token y consultar módulos
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | \
  grep -o '"token":"[^"]*"' | cut -d'"' -f4)

curl http://localhost:8080/api/modules \
  -H "Authorization: Bearer $TOKEN"

# 7. Abrir el frontend
open http://localhost:8080   # macOS
# xdg-open http://localhost:8080  # Linux

# 8. Apagar al terminar
docker compose down
```

---

## Cómo construir y escanear la imagen manualmente

```bash
# Construir
docker build -t TU_USUARIO/curso-gitops:manual .

# Ver el tamaño — debe ser ~20MB
docker images TU_USUARIO/curso-gitops:manual

# Escanear vulnerabilidades (debe dar 0 HIGH/CRITICAL gracias a Alpine + Go)
trivy image --severity HIGH,CRITICAL --format table TU_USUARIO/curso-gitops:manual

# Correr la imagen sola (necesita una MySQL disponible)
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_USER=curso_app \
  -e DB_PASSWORD=C4rs0_S3cur3_P@ss! \
  -e DB_NAME=curso_db \
  TU_USUARIO/curso-gitops:manual
```

---

## Detalles técnicos del Dockerfile (EP10)

**¿Por qué multi-stage?**
Si compilas y ejecutas en la misma imagen `golang:1.22`, el resultado pesa ~1.1GB porque incluye todo el SDK, el compilador y las herramientas. Con multi-stage, el stage final solo contiene el binario compilado y Alpine (~5MB de SO). La imagen resultante pesa ~20MB.

**¿Por qué `CGO_ENABLED=0`?**
CGO es el mecanismo que permite a Go llamar código C. Si está habilitado, el binario depende de `glibc` del sistema operativo donde se compiló. Al deshabilitarlo, el binario es completamente estático — funciona en cualquier Linux, incluyendo Alpine que usa `musl` en lugar de `glibc`.

**¿Por qué correr como usuario no-root?**
Si un atacante explota una vulnerabilidad en la app y logra ejecutar código, que el proceso corra como `appuser` en lugar de `root` limita el daño que puede hacer dentro del contenedor.

**El HEALTHCHECK del Dockerfile:**
```dockerfile
HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1
```
Kubernetes usa los `readinessProbe` y `livenessProbe` del `deployment.yaml` en lugar de este HEALTHCHECK, pero el HEALTHCHECK sirve cuando la imagen corre con `docker run` o `docker compose`.
