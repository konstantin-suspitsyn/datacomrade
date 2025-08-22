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

## Описвние схемы

Описание схемы данных расположено в ```diagrams/data-catalog.drawio```

## Тестирование 

Тесты репозитория строятся с помощью ```testcontainers```. Поднимается тестовый контейнер в модуле ```testcntr```<br>
Файл для инициализации базы данных в папке ```./migrations/```

## Модели данных

## API

### Пользователи
```POST /v1/users/```

#### Регистрация пользователей
```
{
    "email": "email@email.ru",
    "name": "TheName",
    "password": "ThePassword"
}
```

```PUT /v1/users/activate```
```
{
    "token":"tokenfromemail"
}
```
