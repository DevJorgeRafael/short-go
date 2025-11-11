# Go Task Easy List

API REST para gestión de tareas con autenticación JWT y sesiones, construida con Go siguiendo principios de Clean Architecture.

## 🚀 Características

- ✅ Autenticación con JWT (Access + Refresh tokens)
- 🔐 Gestión de sesiones activas
- 📝 CRUD completo de tareas
- 🎯 Sistema de prioridades (Baja, Media, Alta)
- 📊 Estados de tareas (Pendiente, En Progreso, Completada)
- 🏗️ Clean Architecture (Domain, Application, Infrastructure)
- 🗄️ SQLite con GORM
- ✔️ Validación de datos con go-playground/validator

## 📁 Estructura del Proyecto
```
go-task-easy-list/
├── config/                      # Configuración global
│   ├── config.go               # Variables de entorno
│   └── database.go             # Conexión a BD
├── internal/
│   ├── auth/                   # Módulo de autenticación
│   │   ├── application/
│   │   │   └── service/        # Lógica de negocio
│   │   ├── domain/
│   │   │   ├── model/          # Entidades
│   │   │   └── repository/     # Interfaces
│   │   └── infrastructure/
│   │       ├── config/         # Wire/DI
│   │       ├── http/handler/   # Controllers
│   │       └── persistence/    # Implementación repos
│   ├── tasks/                  # Módulo de tareas
│   │   ├── application/
│   │   ├── domain/
│   │   └── infrastructure/
│   └── shared/                 # Código compartido
│       ├── context/            # Context helpers
│       ├── http/               # Response handlers
│       ├── infrastructure/     # Middleware, DI
│       └── validation/         # Validadores
└── migrations/
    └── schema.sql              # Schema de BD
```

## 🛠️ Tecnologías

- **Go 1.23+**
- **Chi** - Router HTTP
- **GORM** - ORM
- **SQLite** - Base de datos
- **JWT** - Autenticación
- **Validator** - Validación de datos

## ⚙️ Instalación

### 1. Clonar el repositorio
```bash
git clone https://github.com/DevJorgeRafael/go-task-easy-list.git
cd go-task-easy-list
```

### 2. Instalar dependencias
```bash
go mod download
```

### 3. Configurar variables de entorno

Copia `.env.example` y configura tus variables:
```bash
cp .env.example .env
```
```env
# Server
PORT=8080

# Database
DB_PATH=./todo.db

# JWT
JWT_SECRET=super-secret-key
JWT_ACCESS_EXPIRATION=1h
JWT_REFRESH_EXPIRATION=7d
```


### 4. Iniciar el servidor
```bash
go run main.go
```

El servidor estará disponible en `http://localhost:8080`

## 📡 API Endpoints

### 🔐 Autenticación (`/api/auth`)

#### Rutas Públicas

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/auth/register` | Registrar nuevo usuario |
| POST | `/api/auth/login` | Iniciar sesión |
| POST | `/api/auth/refresh` | Renovar access token |

#### Rutas Protegidas (requieren JWT)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/auth/logout` | Cerrar sesión |
| GET | `/api/auth/sessions` | Listar sesiones activas |

### ✅ Tareas (`/api/tasks`)

Todas las rutas requieren autenticación (Header: `Authorization: Bearer <token>`)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/tasks` | Crear tarea |
| GET | `/api/tasks` | Listar todas las tareas del usuario |
| GET | `/api/tasks/{id}` | Obtener tarea por ID |
| PUT | `/api/tasks/{id}` | Actualizar tarea |
| DELETE | `/api/tasks/{id}` | Eliminar tarea |
| GET | `/api/tasks/by-status/{statusId}` | Filtrar por estado *(próximamente)* |
| GET | `/api/tasks/by-priority/{priorityId}` | Filtrar por prioridad *(próximamente)* |


## 📊 Modelos de Datos

### Task
```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "statusId": 1,           // 1=Pendiente, 2=En Progreso, 3=Completada
  "priorityId": 2,         // 1=Baja, 2=Media, 3=Alta
  "startsAt": "2025-11-10T09:00:00Z",
  "dueDate": "2025-11-15T18:00:00Z",
  "createdAt": "2025-11-09T22:00:00Z",
  "updatedAt": "2025-11-09T22:00:00Z"
}
```

### Estados (task_statuses)

| ID | Code | Name |
|----|------|------|
| 1 | PENDING | Pendiente |
| 2 | IN_PROGRESS | En Progreso |
| 3 | COMPLETED | Completada |

### Prioridades (task_priorities)

| ID | Code | Name | Level |
|----|------|------|-------|
| 1 | LOW | Baja | 1 |
| 2 | MEDIUM | Media | 2 |
| 3 | HIGH | Alta | 3 |

## 🔒 Seguridad

- Contraseñas hasheadas con bcrypt
- JWT con expiración configurable
- Refresh tokens para renovación segura
- Validación de sesiones activas
- Middleware de autenticación en todas las rutas protegidas


## 👤 Autor

Jorge Rafael Rosero - Proyecto de aprendizaje Go