CREATE SCHEMA dc;
CREATE SCHEMA tech;
CREATE TABLE dc.domain_roles (
    id bigint NOT NULL,
    name character varying(128) NOT NULL,
    description character varying(2000) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL
);
CREATE TABLE dc.table_roles (
    id bigint NOT NULL,
    name character varying(128) NOT NULL,
    description character varying(2000) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL
);
create table dc.user_domain_roles (
                          id bigint primary key not null default nextval('user_domain_roles_id_seq'::regclass),
                          user_id bigint not null,
                          domain_roles_id bigint not null,
                          created_at timestamp without time zone not null default now(),
                          updated_at timestamp without time zone not null default now(),
                          is_deleted boolean not null default false,
                          domain_id bigint not null,
                          updated_by_id bigint not null
);
create table dc.user_table_roles
(
    id             bigserial not null,
    user_id        bigint                  not null,
    table_roles_id bigint                  not null,
    created_at     timestamp default now() not null,
    updated_at     timestamp default now() not null,
    is_deleted     boolean   default false not null,
    table_id       bigint                  not null,
    updated_by_id  bigint                  not null
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

