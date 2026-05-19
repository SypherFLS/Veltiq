internal/
├── core/                 # ядро приложения: module system, pipeline, lifecycle
│   ├── module/
│   ├── pipeline/
│   ├── contracts/
│   └── lifecycle/
│
├── infrastructure/       # технические реализации
│   ├── config/
│   │   └── loader.go
│   │
│   ├── logging/
│   │   ├── logger.go
│   │   └── context.go
│   │
│   ├── database/
│   │   ├── postgres/
│   │   │   ├── connection.go
│   │   │   └── migrations.go
│   │   └── transaction.go
│   │
│   ├── cache/
│   │   └── redis.go
│   │
│   ├── http/
│   │   ├── server.go
│   │   ├── router.go
│   │   └── middleware/
│   │       ├── logging.go
│   │       ├── recovery.go
│   │       ├── cors.go
│   │       └── auth.go
│   │
│   ├── jwt/
│   │   └── manager.go
│   │
│   ├── cookies/
│   │   └── manager.go
│   │
│   └── security/
│       ├── password.go
│       └── random.go
│
├── modules/
│   ├── auth/
│   ├── users/
│   └── billing/
│
└── app/
    └── bootstrap.go