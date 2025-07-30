# Документация для разработчика

## Структура
├─cmd/
├─configs/
├─data/
│ ├─dataerrors/
│ ├─domainmodel/
│ ├─domainrolemodel/
│ ├─rolesmodel/
│ ├─rolesmodel/
├─db/
├─diagrams/
├─internal/
│ ├─healthcheck/
│ ├─services/
│ ├─users/
│ ├─utils/
│ │ ├─ custresponse/
│ │ ├─ jsonlog/
│ │ ├─ mailer/
│ │ ├─ shared/
│ │ └─ validator/
└─...

## Обязательные файлы для запуска

```.dev-env``` нужен для правильной работы testcontainers
Чтобы узнать какие переменные нужны, смотри файл ```.env-example```

## Модели данных

