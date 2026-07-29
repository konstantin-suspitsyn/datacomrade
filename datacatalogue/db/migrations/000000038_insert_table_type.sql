-- +goose up
INSERT INTO dc.table_type
("name", description, is_deleted, created_at, updated_at, user_id)
values
('fact', 'Фактовая таблица', false, now(), now(), 1),
('dimension', 'Словарь', false, now(), now(), 1);
