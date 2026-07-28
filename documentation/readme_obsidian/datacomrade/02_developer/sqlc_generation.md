2026-03-30 20:46

Status:

Tags:  [[new_microservice]] [[db]]

# sqlc_generation

Если вам нужны миграции для sqlc, на данный момент из Docker Postgres, то можно воспользоваться ```deploy/compose/sqlc/generate_sqlc_files.py``` для формирования cхемы для sqlc и yml файлов. Настройка находится в ```deploy/compose/sqlc/schema_to_create.yml```

```
python generate_sqlc_files.py schema_to_create.yml ../datacatalogue/.env
```

Структура ```schema_to_create.yml```

```yaml
название микросервиса:
  save_path: "internal/repository" #путь куда сохранять код, который сформирует sqlc (лучше не менять)
  env_path: "путь env файлу микросервиса"
  models:
    название_модели: # будет использована как папка для schema.sql
      name: имя_модели # будет использована при генерации кода
      tables:
        - схема.таблица # перечисляем все таблицы, которые будут внутри класса

```


Запускать следующим образом ```python generate_sqlc_files.py schema_to_create.yml путь_к_env```

В результате будут сформированы файлы:
```
.
├── микросервис
│   ├── db
│   │   ├── sqlc
│   │   │   └── название модели
│   │   │   │   └── schema.sql
│   │   │   └── название модели 1
│   │   │   │   └── schema.sql
│   │   └── sqlc.yml

```
### References
