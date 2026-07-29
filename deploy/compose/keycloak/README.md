# Keycloak

Сервер авторизации и аутентификации для datacomrade: Keycloak 26 + собственный Postgres.

## Запуск

Стек подключается к внешней сети `microservices-net`, которую создаёт core-стек, поэтому core должен быть поднят первым.

Проще всего поднять всё сразу:

```bash
task up-all
```

Либо по шагам — bash:

```bash
task up-core && task up-keycloak
```

PowerShell (в 5.1 нет оператора `&&`):

```powershell
task up-core; if ($?) { task up-keycloak }
```

Остановка — `task down-keycloak` (данные остаются в volume `keycloak_kc_pgdata`).

## Что где

| Что | Адрес |
| --- | --- |
| Админка | http://localhost:8081 |
| OIDC discovery | http://localhost:8081/realms/{realm}/.well-known/openid-configuration |
| Health | http://localhost:9000/health/ready |
| Metrics | http://localhost:9000/metrics |
| Postgres | localhost:5433 |

Логин в админку — `KEYCLOAK_ADMIN_USER` / `KEYCLOAK_ADMIN_PASSWORD` из `.env` (по умолчанию `admin` / `admin`).

## Подключение Go-сервисов

Из контейнера в сети `microservices-net` Keycloak доступен по DNS-имени `keycloak`:

```
issuer:  http://keycloak:8080/realms/datacomrade
jwks:    http://keycloak:8080/realms/datacomrade/protocol/openid-connect/certs
```

С хоста (при локальном запуске сервиса вне докера) — те же пути, но `http://localhost:8081`.

Обратите внимание: issuer в токене будет содержать тот хост, по которому шёл запрос. Если сервис валидирует токены, полученные через `localhost:8081`, а сам ходит в `keycloak:8080`, issuer не совпадёт. Для смешанного режима проще выставить `KEYCLOAK_HOSTNAME` и обращаться к Keycloak по одному и тому же адресу отовсюду.

## Важные оговорки

- **Режим `start-dev`.** Только для локальной разработки: HTTP без TLS, кеши в памяти, отключена строгая проверка hostname. Для прода нужен `start --optimized` с отдельным Dockerfile (`kc.sh build`), TLS и `KC_HOSTNAME_STRICT=true`.
- **Админ создаётся один раз.** `KC_BOOTSTRAP_ADMIN_*` срабатывают только на пустой БД. Чтобы пересоздать — `docker compose down -v` (удалит все realm'ы и пользователей).
- **Первый старт долгий.** Keycloak накатывает свои миграции, поэтому у healthcheck `start_period: 60s`.
- **Realm нужно создать вручную** через админку или `kcadm.sh`. Если понадобится воспроизводимая конфигурация — realm можно экспортировать в `realm-export.json` и импортировать при старте через `--import-realm`.
