CREATE SCHEMA dc;
CREATE SCHEMA tech;
CREATE TABLE dc."user" (
    id bigint NOT NULL,
    name character varying(512) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    incoming_user_id bigint DEFAULT '-1'::integer NOT NULL
);
