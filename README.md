subscriptions-service/
│
├─ cmd/
│   └─ app/
│       └─ main.go              # точка входа
│
├─ internal/
│   ├─ config/
│   │   └─ config.go            # загрузка .env / yaml
│   │
│   ├─ handler/                 # HTTP слой
│   │   ├─ subscription.go
│   │   ├─ report.go            # сумма подписок
│   │   └─ router.go
│   │
│   ├─ service/                 # бизнес-логика
│   │   ├─ subscription.go
│   │   └─ report.go
│   │
│   ├─ repository/              # работа с БД
│   │   ├─ subscription_pg.go
│   │   └─ repository.go
│   │
│   ├─ model/                   # структуры данных
│   │   └─ subscription.go
│   │
│   ├─ logger/
│   │   └─ logger.go
│   │
│   └─ middleware/
│       └─ logging.go
│
├─ migrations/
│   ├─ 000001_init.up.sql
│   └─ 000001_init.down.sql
│
├─ docs/                        # swagger
│   ├─ docs.go
│   ├─ swagger.json
│   └─ swagger.yaml
│
├─ pkg/                         # переиспользуемые утилиты (опционально)
│   └─ postgres/
│       └─ postgres.go
│
├─ .env
├─ docker-compose.yml
├─ Dockerfile
├─ Makefile
├─ go.mod
└─ README.md# sub-test
