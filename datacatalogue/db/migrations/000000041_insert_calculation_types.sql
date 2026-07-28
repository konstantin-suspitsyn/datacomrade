-- +goose up
INSERT INTO dc.calculation_type
("name", description, created_at, updated_at, is_deleted)
values
('sum', 'Сумма', now(), now(), false),
('count', 'Количество', now(), now(), false),
('count_distinct', 'Количество уникальных', now(), now(), false),
('avg', 'Среднее', now(), now(), false),
('max', 'Максимум', now(), now(), false),
('min', 'Минимум', now(), now(), false)
;
