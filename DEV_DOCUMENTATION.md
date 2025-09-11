# Документация для разработчика

## Структура
```
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
```

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

## Ролевая модель авторизации
На данный момент, модель авторизация будет достаточно простая и состоит из действий, которые может совершать пользователь. Данные находятся в таблице ```users.action```.

Также роли должны быть продублированы в пакете ```config```

1. write - создание записей в БД
2. read - чтение записей в БД
3. update - обновление записей
4. deleteOwn - удаление записей
5. readOwn - чтение собственных записей в БД
6. updateOwn - обновление собственных записей
7. deleteOwn - удаление собственных записей

Доступ осуществляется либо на уровне домена (если он есть), либо на уровне конкретных таблиц

Модуль авторизации проверит первичное наличие доступа к записи, и либо пропустит к select, либо выдаст ошибку

Предполагается, что основной доступ будет осуществляться через роли. Роль получит разрешение на действие к определенному объекту по домему или по id. Для этого используется таблица ```users.role_access```. Поле ```domain_id```. Если domain_id == 1, значит, доступ не по домену, а по объекту

Также будет возможность выдать роль конкретному пользователю. Для этого используется таблица ```users.user_access```, однако, доступ через эту таблицу будет осуществляться только к записям таблиц, без доменов. Данный способ не приветствуется.

Лучший способ - роль с привязанным доменом



## Заметки для формирования БД
Эти вещи должны быть вставлены после создания структуры БД
```sql
-- Создаем базовые роли для авторизации
INSERT INTO users."action" ("name",description,is_deleted,created_at,updated_at) VALUES
	 ('write','Write data to DB',false,now(),now()),
	 ('read','Read data from DB',false,now(), now()),
	 ('update','Update data in DB',false,now(), now()),
	 ('delete','Mark data deleted',false,now(), now()),
	 ('readOwn','Read your own data',false,now(), now()),
	 ('updateOwn','Update own data',false,now(), now()),
	 ('deleteOwn','Delete own data',false,now(), now()),
	 ('shareOwn','Share resource with others',false,now(), now()),
	 ('query','Query data',false,now(), now());

```

```sql


INSERT INTO shared."domain" ("name",description,is_deleted,created_at,updated_at) VALUES
	 ('None','Placeholder for none value',false,now(), now());

```
