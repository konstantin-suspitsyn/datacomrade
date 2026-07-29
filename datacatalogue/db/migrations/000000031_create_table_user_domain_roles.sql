-- +goose up
create table dc.user_domain_roles
(
    id              bigserial
        constraint user_domain_roles_pk
            primary key,
    user_id         bigint                  not null
        constraint user_domain_roles_user_fk
            references dc."user"
            on update cascade on delete restrict,
    domain_roles_id bigint                  not null
        constraint user_domain_roles_domain_roles_fk
            references dc.domain_roles,
    created_at      timestamp default now() not null,
    updated_at      timestamp default now() not null,
    is_deleted      boolean   default false not null,
    domain_id       bigint                  not null
        constraint user_domain_roles_domain_cat_id_fk
            references dc.domain_cat
);

comment on constraint user_domain_roles_domain_cat_id_fk on dc.user_domain_roles is 'domain id fk';

alter table dc.user_domain_roles
    owner to postgres_user;


-- +goose down
DROP TABLE IF EXISTS dc.user_domain_roles; 