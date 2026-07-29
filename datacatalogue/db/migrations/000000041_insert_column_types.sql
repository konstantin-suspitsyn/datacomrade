-- +goose up
INSERT INTO dc.column_type
("name", description, is_deleted, created_at, updated_at, user_id)
SELECT v.name, v.description, false, now(), now(), u.id
FROM (VALUES ('int', 'Целое 32 бита'),
             ('bigint', 'Целое 64 бита, идентификаторы'),
             ('decimal', 'Точное десятичное, деньги'),
             ('float', 'Приблизительное с плавающей точкой'),
             ('bool', 'Логический флаг'),
             ('string', 'Короткая строка, годится для группировки'),
             ('text', 'Длинный текст, не для группировки'),
             ('date', 'Календарная дата без времени'),
             ('datetime', 'Дата и время без часового пояса'),
             ('datetime_tz', 'Дата и время с часовым поясом')
     ) AS v(name, description)
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'Lomonosov M.') u;

-- +goose down
DELETE FROM dc.column_type
WHERE "name" IN ('int', 'bigint', 'decimal', 'float', 'bool',
                 'string', 'text', 'date', 'datetime', 'datetime_tz');
