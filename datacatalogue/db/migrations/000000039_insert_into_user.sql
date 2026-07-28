-- +goose up
INSERT INTO dc."user"
("name", created_at, updated_at, is_deleted)
VALUES('Lomonosov M.', now(), now(), false);

-- +goose down
DELETE FROM dc."user"
WHERE "name" = 'Lomonosov M.'
