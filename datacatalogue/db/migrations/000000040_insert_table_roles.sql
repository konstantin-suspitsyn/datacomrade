-- +goose up
INSERT INTO dc.table_roles (id, name, description, created_at, updated_at, is_deleted)
VALUES (DEFAULT, 'can_read'::varchar(128), 'Пользователь может читать все поля и таблицы из домена'::varchar(2000),
        DEFAULT, DEFAULT, DEFAULT);

INSERT INTO dc.table_roles (id, name, description, created_at, updated_at, is_deleted)
VALUES (DEFAULT, 'can_write'::varchar(128), 'Пользователь может изменять таблицы и поля домена'::varchar(2000), DEFAULT,
        DEFAULT, DEFAULT);

INSERT INTO dc.table_roles (id, name, description, created_at, updated_at, is_deleted)
VALUES (DEFAULT, 'can_grant'::varchar(128), 'Пользователь может выдавать права на домен'::varchar(2000), DEFAULT,
        DEFAULT, DEFAULT);

-- +goose down
DELETE FROM dc.table_roles
WHERE "name" in ('can_read', 'can_write', 'can_grant')
