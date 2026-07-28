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
CREATE TABLE dc.domains_domain_roles (
    id bigint NOT NULL,
    domain_cat_id bigint NOT NULL,
    domain_roles_id bigint NOT NULL,
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
CREATE TABLE dc.tables_table_roles (
    id bigint NOT NULL,
    table_cat_id bigint NOT NULL,
    table_roles_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL
);
CREATE TABLE dc.user_domain_roles (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    domain_roles_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL
);
CREATE TABLE dc.user_table_roles (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    table_roles_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL
);
