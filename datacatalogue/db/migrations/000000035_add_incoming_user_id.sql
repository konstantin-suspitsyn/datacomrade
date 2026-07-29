-- +goose up
ALTER TABLE dc."user" ADD incoming_user_id bigint DEFAULT -1 NOT NULL;
COMMENT ON COLUMN dc."user".incoming_user_id IS 'User id from user microservice';
ALTER TABLE dc."user" ADD CONSTRAINT incoming_user_id_unique UNIQUE (incoming_user_id);

