-- +goose up
create table dc.user_table_roles
(
    id             bigserial
        constraint user_table_roles_pk
            primary key,
    user_id        bigint                  not null
        constraint user_table_roles_user_fk
            references dc."user"
            on update cascade on delete restrict,
    table_roles_id bigint                  not null
        constraint user_table_roles_table_roles_fk
            references dc.table_roles
            on update cascade on delete restrict,
    created_at     timestamp default now() not null,
    updated_at     timestamp default now() not null,
    is_deleted     boolean   default false not null,
    table_id       bigint                  not null
        constraint user_table_roles_table_cat_id_fk
            references dc.table_cat
);

-- +goose down
DROP TABLE IF EXISTS dc.user_table_roles; 