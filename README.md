# Short Go

API REST para acortar enlaces, generar códigos QR y visualizar analíticas, construida con Go siguiendo principios de Clean Architecture y Arquitectura Modular.

## 🚀 Características

- ✅ Autenticación con JWT (Access + Refresh tokens)
- 🔐 Gestión de sesiones activas y recuperación de contraseña vía Email
- 🔗 Acortador de URLs con redirección eficiente
- 📊 Sistema de analíticas y rastreo de clicks
- 📱 Generación de códigos QR dinámicos
- 🏗️ Arquitectura Modular (Auth, ShortLinks, Analytics, QR)
- 🗄️ PostgreSQL con GORM
- ✔️ Validación de datos con go-playground/validator

## 📁 Estructura del Proyecto
```
short-go/
├── config/                      # Configuración global
│   ├── config.go               # Variables de entorno
│   └── database.go             # Conexión a BD y GORM AutoMigrate
├── internal/
│   ├── analytics/               # Módulo de analíticas
│   │   ├── application/
│   │   │   └── service/        # Lógica de registro de clicks
│   │   ├── domain/
│   │   │   ├── model/          # Entidades (Click)
│   │   │   └── repository/     # Interfaces
│   │   └── infrastructure/
│   │       ├── config/         # Wire/DI del módulo
│   │       ├── http/handler/   # Controllers
│   │       └── persistence/    # Implementación GORM
│   ├── auth/                    # Módulo de autenticación
│   │   ├── application/
│   │   │   └── service/        # Lógica de login/register
│   │   ├── domain/
│   │   │   ├── model/          # Entidades (User, Session)
│   │   │   └── repository/     # Interfaces
│   │   └── infrastructure/
│   │       ├── config/         # Wire/DI del módulo
│   │       ├── email/          # Servicio de envío (Brevo)
│   │       ├── http/handler/   # Controllers
│   │       └── persistence/    # Implementación GORM
│   ├── qr/                      # Módulo de códigos QR
│   │   └── infrastructure/
│   │       ├── config/         # Wire/DI del módulo
│   │       └── http/handler/   # Generador de imágenes QR
│   ├── short-links/             # Módulo de acortador
│   │   ├── application/
│   │   │   └── service/        # Lógica de creación/redirección
│   │   ├── domain/
│   │   │   ├── model/          # Entidades (ShortLink)
│   │   │   └── repository/     # Interfaces
│   │   └── infrastructure/
│   │       ├── config/         # Wire/DI del módulo
│   │       ├── http/handler/   # Controllers
│   │       └── persistence/    # Implementación GORM
│   └── shared/                  # Código compartido
│       ├── context/            # Context helpers
│       ├── http/               # Response helpers
│       ├── infrastructure/     # Container DI y Middleware
│       └── validation/         # Validadores personalizados
├── .env                        # Variables de entorno (local)
├── .env.template               # Plantilla de variables
├── go.mod                      # Dependencias
└── main.go                     # Entry point
```

## 🛠️ Tecnologías

- **Go 1.25+**
- **Chi v5** - Router HTTP ligero y rápido
- **GORM** - ORM robusto para Go
- **PostgreSQL** - Base de datos relacional
- **JWT v5** - Autenticación y seguridad
- **Validator v10** - Validación de datos y estructuras
- **Go QR Code** - Generación de códigos QR nativa
- **UUID** - Generación de identificadores únicos
- **Bcrypt** - Hashing seguro de contraseñas
- **Godotenv** - Carga de variables de entorno

## ⚙️ Instalación

### 1. Clonar el repositorio
```bash
git clone [https://github.com/DevJorgeRafael/short-go.git](https://github.com/DevJorgeRafael/short-go.git)
cd short-go
```

### 2. Instalar dependencias
```bash
go mod download
```

### 3. Configurar variables de entorno

Copia `.env.example` y configura las variables de entorno:
```bash
cp .env.example .env
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
| POST | `/api/auth/login` | Iniciar sesión (Retorna Access + Refresh Token) |
| POST | `/api/auth/refresh` | Renovar Access Token |
| POST | `/api/auth/forgot-password` | Solicitar correo de recuperación de contraseña |
| POST | `/api/auth/reset-password` | Restablecer contraseña usando token |

#### Rutas Protegidas (requieren JWT)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/auth/logout` | Cerrar sesión actual |
| GET | `/api/auth/sessions` | Listar sesiones activas del usuario |

### 🔗 Short Links

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/short-links` | Crear enlace corto (Auth opcional para asociar al usuario) |
| GET | `/{code}` | Redireccionar a la URL original (Ruta Raíz) |

### 📊 Analíticas (`/api/stats`)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/stats/{code}` | Obtener estadísticas y contador de clicks |

### 📱 Códigos QR (`/api/qr`)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/qr/{code}` | Generar imagen del código QR para un enlace |


## 🔒 Seguridad

- **Contraseñas Seguras**: Hasheadas con bcrypt antes de ser almacenadas.
- **Autenticación JWT**: Implementación de Access Tokens y Refresh Tokens con tiempos de expiración configurables.
- **Gestión de Sesiones**: Control y validación de sesiones activas en base de datos.
- **Recuperación de Contraseña**: Envío de códigos vía Email (Brevo API). Por seguridad, los códigos de verificación se guardan hasheados en la base de datos, nunca en texto plano.
- **Middleware de Protección**: Verificación de autenticación en todas las rutas protegidas.


## 👤 Autor

Jorge Rafael Rosero - Acortador de enlaces con Go