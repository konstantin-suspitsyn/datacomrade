-- +goose up

alter table dc."user"
drop column incoming_user_id;