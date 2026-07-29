-- +goose up
ALTER TABLE dc."user" ADD external_id uuid NOT NULL;
COMMENT ON COLUMN dc."user".external_id IS 'Subject (sub) from Keycloak';
ALTER TABLE dc."user" ADD CONSTRAINT user_external_id_unique UNIQUE (external_id);
CREATE INDEX user_external_id_idx ON dc."user" (external_id);

-- +goose down
DROP INDEX IF EXISTS dc.user_external_id_idx;
ALTER TABLE dc."user" DROP CONSTRAINT IF EXISTS user_external_id_unique;
ALTER TABLE dc."user" DROP COLUMN IF EXISTS external_id;
