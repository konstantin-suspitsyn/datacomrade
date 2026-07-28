2026-03-18 21:50

Status:

Tags: [[new_microservice]] [[env]]

# Работа с env файлами

Если вы создаете новый микросервис с настройками в .env, то работа следующая:
```
.
├── deploy
│   ├── compose
│   │   ├── core
│   │   │   └── docker-compose.yml
│   │   └── datacatalogue
│   │       └── docker-compose.yml
│   └── env
│       ├── datacatalogue.env.template
│       ├── datacatalogue.env
│       ├── .env
│       ├── .env.template
│       └── generate-env.sh
```

.env управляются централизовано. Это файл .env и .env.template. В них должны быть все переменные

.env
```yml

# ----------------------------------------
# Данные для микросервиса DATACATALOGUE
# ----------------------------------------

# ----------------------------------------
# GRPC host and port
# ----------------------------------------
DATACATALOGUE_GRPC_HOST=localhost
DATACATALOGUE_GRPC_PORT=5051

# ----------------------------------------
# Logger settings
# ----------------------------------------
DATACATALOGUE_LOGGER_LEVEL=info
DATACATALOGUE_LOGGER_AS_JSON=true

# ----------------------------------------
# Postgres DB settings
# ----------------------------------------
DATACATALOGUE_POSTGRES_IMAGE=18.3-alpine3.23
DATACATALOGUE_POSTGRES_USER=postgres_user
DATACATALOGUE_POSTGRES_PASSWORD=pass
DATACATALOGUE_POSTGRES_HOST=localhost
# Внешний порт Postgres
DATACATALOGUE_POSTGRES_PORT=5432

DATACATALOGUE_POSTGRES_DB=dcpostgres

# ----------------------------------------
# Данные для других микросервисов
# ----------------------------------------
```

.env.template (файл для примера людям, чтобы они знали как заполнять .env)
```yml

# ----------------------------------------
# Данные для микросервиса DATACATALOGUE
# ----------------------------------------

# ----------------------------------------
# GRPC host and port
# ----------------------------------------
DATACATALOGUE_GRPC_HOST=localhost
DATACATALOGUE_GRPC_PORT=5051

# ----------------------------------------
# Logger settings
# ----------------------------------------
DATACATALOGUE_LOGGER_LEVEL=info
DATACATALOGUE_LOGGER_AS_JSON=true

# ----------------------------------------
# Postgres DB settings
# ----------------------------------------
DATACATALOGUE_POSTGRES_IMAGE=18.3-alpine3.23
DATACATALOGUE_POSTGRES_USER=postgres_user
DATACATALOGUE_POSTGRES_PASSWORD=pass
DATACATALOGUE_POSTGRES_HOST=localhost
# Внешний порт Postgres
DATACATALOGUE_POSTGRES_PORT=5432

DATACATALOGUE_POSTGRES_DB=dcpostgres

# ----------------------------------------
# Данные для других микросервисов
# ----------------------------------------
```

```datacatalogue.env``` - шаблон для переноса .env файл в ```deploy/compose/datacatalogue/.env```
То есть ```{имя микросервиса}.env``` - шаблон для переноса .env файл в ```deploy/compose/{имя микросервиса}/.env```

```datacatalogue.env.template``` - шаблон для переноса .env файл в ```deploy/compose/datacatalogue/.env.template```
То есть ```{имя микросервиса}.env.template``` - шаблон для переноса .env файл в ```deploy/compose/{имя микросервиса}/.env.template```

Сам шаблон выглядит так:
```yml

# ----------------------------------------
# Данные для микросервиса DATACATALOGUE
# ----------------------------------------

# ----------------------------------------
# GRPC host and port
# ----------------------------------------
DATACATALOGUE_GRPC_HOST=${DATACATALOGUE_GRPC_HOST}
DATACATALOGUE_GRPC_PORT=${DATACATALOGUE_GRPC_PORT}

# ----------------------------------------
# Logger settings
# ----------------------------------------
DATACATALOGUE_LOGGER_LEVEL=${DATACATALOGUE_LOGGER_LEVEL}
DATACATALOGUE_LOGGER_AS_JSON=${DATACATALOGUE_LOGGER_AS_JSON}

# ----------------------------------------
# Postgres DB settings
# ----------------------------------------
DATACATALOGUE_POSTGRES_IMAGE=${DATACATALOGUE_POSTGRES_IMAGE}
DATACATALOGUE_POSTGRES_USER=${DATACATALOGUE_POSTGRES_USER}
DATACATALOGUE_POSTGRES_PASSWORD=${DATACATALOGUE_POSTGRES_PASSWORD}
DATACATALOGUE_POSTGRES_HOST=${DATACATALOGUE_POSTGRES_HOST}
# Внешний порт Postgres
DATACATALOGUE_POSTGRES_PORT=${DATACATALOGUE_POSTGRES_PORT}

DATACATALOGUE_POSTGRES_DB=${DATACATALOGUE_POSTGRES_DB}
```
${название переменной} берется либо из .env или .env.template и создает файлы внутри микросервисов

Реализуется автоматизация через ```generate-env.sh```
Но запускать нужно через [[Taskdetail]]. Команда ```task env:generate```
### References