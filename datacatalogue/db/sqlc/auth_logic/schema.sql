CREATE SCHEMA dc;
CREATE SCHEMA tech;

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

create table dc.domain_roles
(
    id          bigserial
        constraint domain_roles_pk
            primary key,
    name        varchar(128)            not null
        constraint domain_roles_name_unique
            unique,
    description varchar(2000)           not null,
    created_at  timestamp default now() not null,
    updated_at  timestamp default now() not null,
    is_deleted  boolean   default false not null
);


create table dc."user"
(
    id          bigserial
        constraint user_pk
            primary key,
    name        varchar(512)            not null
        constraint user_name_unique
            unique,
    created_at  timestamp default now() not null,
    updated_at  timestamp default now() not null,
    is_deleted  boolean   default false not null,
    external_id uuid                    not null
        constraint user_external_id_unique
            unique
);


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


create table dc.table_roles
(
    id          bigserial
        constraint table_roles_pk
            primary key,
    name        varchar(128)            not null
        constraint table_roles_name_unique
            unique,
    description varchar(2000)           not null,
    created_at  timestamp default now() not null,
    updated_at  timestamp default now() not null,
    is_deleted  boolean   default false not null
);

create table dc.table_cat
(
    id            bigserial
        constraint table_cat_pk
            primary key,
    name          varchar(128)            not null,
    description   varchar(2000)           not null,
    schema_id     bigint                  not null
        constraint table_cat_schema_cat_fk
            references dc.schema_cat
            on update cascade on delete restrict,
    table_type_id bigint                  not null
        constraint table_cat_table_type_fk
            references dc.table_type
            on update cascade on delete restrict,
    domain_id     bigint                  not null
        constraint table_cat_domain_cat_fk
            references dc.domain_cat
            on update cascade on delete restrict,
    is_deleted    boolean   default false not null,
    is_get_dict   boolean   default false not null,
    created_at    timestamp default now() not null,
    updated_at    timestamp default now() not null,
    user_id       bigint                  not null
        constraint table_cat_user_fk
            references dc."user"
            on update cascade on delete restrict,
    constraint table_cat_name_unique
        unique (schema_id, name)
);

