-- +goose up
INSERT INTO dc."user"
("name", created_at, updated_at, is_deleted, external_id)
VALUES('Lomonosov M.', now(), now(), false, '00000000-0000-0000-0000-000000000000');

-- +goose down
DELETE FROM dc."user"
WHERE "name" = 'Lomonosov M.'
