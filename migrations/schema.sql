--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4 (Debian 17.4-1.pgdg120+2)
-- Dumped by pg_dump version 17.4 (Debian 17.4-1.pgdg120+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: datacatalog; Type: SCHEMA; Schema: -; Owner: postgres_user
--

CREATE SCHEMA datacatalog;


ALTER SCHEMA datacatalog OWNER TO postgres_user;

--
-- Name: users; Type: SCHEMA; Schema: -; Owner: postgres_user
--

CREATE SCHEMA users;


ALTER SCHEMA users OWNER TO postgres_user;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: alias_catalogue; Type: TABLE; Schema: datacatalog; Owner: postgres_user
--

CREATE TABLE datacatalog.alias_catalogue (
    id bigint NOT NULL,
    name character varying NOT NULL,
    description character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at time with time zone DEFAULT now() NOT NULL,
    updated_at time with time zone DEFAULT now() NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE datacatalog.alias_catalogue OWNER TO postgres_user;

--
-- Name: COLUMN alias_catalogue.name; Type: COMMENT; Schema: datacatalog; Owner: postgres_user
--

COMMENT ON COLUMN datacatalog.alias_catalogue.name IS 'This alias will be used to get all tables for user';


--
-- Name: alias_id_seq; Type: SEQUENCE; Schema: datacatalog; Owner: postgres_user
--

CREATE SEQUENCE datacatalog.alias_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE datacatalog.alias_id_seq OWNER TO postgres_user;

--
-- Name: alias_id_seq; Type: SEQUENCE OWNED BY; Schema: datacatalog; Owner: postgres_user
--

ALTER SEQUENCE datacatalog.alias_id_seq OWNED BY datacatalog.alias_catalogue.id;


--
-- Name: column_catalogue; Type: TABLE; Schema: datacatalog; Owner: postgres_user
--

CREATE TABLE datacatalog.column_catalogue (
    id bigint NOT NULL,
    name character varying NOT NULL,
    description character varying NOT NULL,
    alias_id bigint NOT NULL,
    table_id bigint NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at time with time zone DEFAULT now() NOT NULL,
    updated_at time with time zone DEFAULT now() NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE datacatalog.column_catalogue OWNER TO postgres_user;

--
-- Name: column_catalogue_id_seq; Type: SEQUENCE; Schema: datacatalog; Owner: postgres_user
--

CREATE SEQUENCE datacatalog.column_catalogue_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE datacatalog.column_catalogue_id_seq OWNER TO postgres_user;

--
-- Name: column_catalogue_id_seq; Type: SEQUENCE OWNED BY; Schema: datacatalog; Owner: postgres_user
--

ALTER SEQUENCE datacatalog.column_catalogue_id_seq OWNED BY datacatalog.column_catalogue.id;


--
-- Name: database_catalogue; Type: TABLE; Schema: datacatalog; Owner: postgres_user
--

CREATE TABLE datacatalog.database_catalogue (
    id bigint NOT NULL,
    name character varying NOT NULL,
    host_id bigint NOT NULL,
    description character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at time with time zone DEFAULT now() NOT NULL,
    updated_at time with time zone DEFAULT now() NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE datacatalog.database_catalogue OWNER TO postgres_user;

--
-- Name: database_id_seq; Type: SEQUENCE; Schema: datacatalog; Owner: postgres_user
--

CREATE SEQUENCE datacatalog.database_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE datacatalog.database_id_seq OWNER TO postgres_user;

--
-- Name: database_id_seq; Type: SEQUENCE OWNED BY; Schema: datacatalog; Owner: postgres_user
--

ALTER SEQUENCE datacatalog.database_id_seq OWNED BY datacatalog.database_catalogue.id;


--
-- Name: database_type_catalogue; Type: TABLE; Schema: datacatalog; Owner: postgres_user
--

CREATE TABLE datacatalog.database_type_catalogue (
    id bigint NOT NULL,
    name character varying NOT NULL,
    db_version character varying NOT NULL,
    description character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at time with time zone DEFAULT now() NOT NULL,
    updated_at time with time zone DEFAULT now() NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE datacatalog.database_type_catalogue OWNER TO postgres_user;

--
-- Name: COLUMN database_type_catalogue.user_id; Type: COMMENT; Schema: datacatalog; Owner: postgres_user
--

COMMENT ON COLUMN datacatalog.database_type_catalogue.user_id IS 'last updated by';


--
-- Name: database_type_id_seq; Type: SEQUENCE; Schema: datacatalog; Owner: postgres_user
--

CREATE SEQUENCE datacatalog.database_type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE datacatalog.database_type_id_seq OWNER TO postgres_user;

--
-- Name: database_type_id_seq; Type: SEQUENCE OWNED BY; Schema: datacatalog; Owner: postgres_user
--

ALTER SEQUENCE datacatalog.database_type_id_seq OWNED BY datacatalog.database_type_catalogue.id;


--
-- Name: domain_catalogue; Type: TABLE; Schema: datacatalog; Owner: postgres_user
--

CREATE TABLE datacatalog.domain_catalogue (
    id bigint NOT NULL,
    name character varying NOT NULL,
    description character varying,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE datacatalog.domain_catalogue OWNER TO postgres_user;

--
-- Name: domain_id_seq; Type: SEQUENCE; Schema: datacatalog; Owner: postgres_user
--

CREATE SEQUENCE datacatalog.domain_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE datacatalog.domain_id_seq OWNER TO postgres_user;

--
-- Name: domain_id_seq; Type: SEQUENCE OWNED BY; Schema: datacatalog; Owner: postgres_user
--

ALTER SEQUENCE datacatalog.domain_id_seq OWNED BY datacatalog.domain_catalogue.id;


--
-- Name: group_levels_catalogue; Type: TABLE; Schema: datacatalog; Owner: postgres_user
--

CREATE TABLE datacatalog.group_levels_catalogue (
    id bigint NOT NULL,
    column_id bigint NOT NULL,
    parent_column_id bigint NOT NULL,
    level smallint NOT NULL,
    description character varying,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at time with time zone DEFAULT now() NOT NULL,
    updated_at time with time zone DEFAULT now() NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE datacatalog.group_levels_catalogue OWNER TO postgres_user;

--
-- Name: group_levels_catalogue_id_seq; Type: SEQUENCE; Schema: datacatalog; Owner: postgres_user
--

CREATE SEQUENCE datacatalog.group_levels_catalogue_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE datacatalog.group_levels_catalogue_id_seq OWNER TO postgres_user;

--
-- Name: group_levels_catalogue_id_seq; Type: SEQUENCE OWNED BY; Schema: datacatalog; Owner: postgres_user
--

ALTER SEQUENCE datacatalog.group_levels_catalogue_id_seq OWNED BY datacatalog.group_levels_catalogue.id;


--
-- Name: host_catalogue; Type: TABLE; Schema: datacatalog; Owner: postgres_user
--

CREATE TABLE datacatalog.host_catalogue (
    id bigint NOT NULL,
    name character varying NOT NULL,
    host_address_env character varying NOT NULL,
    port character varying NOT NULL,
    database_type_id bigint NOT NULL,
    description character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at time with time zone DEFAULT now() NOT NULL,
    updated_at time with time zone DEFAULT now() NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE datacatalog.host_catalogue OWNER TO postgres_user;

--
-- Name: host_id_seq; Type: SEQUENCE; Schema: datacatalog; Owner: postgres_user
--

CREATE SEQUENCE datacatalog.host_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE datacatalog.host_id_seq OWNER TO postgres_user;

--
-- Name: host_id_seq; Type: SEQUENCE OWNED BY; Schema: datacatalog; Owner: postgres_user
--

ALTER SEQUENCE datacatalog.host_id_seq OWNED BY datacatalog.host_catalogue.id;


--
-- Name: schema_catalogue; Type: TABLE; Schema: datacatalog; Owner: postgres_user
--

CREATE TABLE datacatalog.schema_catalogue (
    id bigint NOT NULL,
    name character varying NOT NULL,
    database_id bigint NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at time with time zone DEFAULT now() NOT NULL,
    updated_at time with time zone DEFAULT now() NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE datacatalog.schema_catalogue OWNER TO postgres_user;

--
-- Name: schema_id_seq; Type: SEQUENCE; Schema: datacatalog; Owner: postgres_user
--

CREATE SEQUENCE datacatalog.schema_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE datacatalog.schema_id_seq OWNER TO postgres_user;

--
-- Name: schema_id_seq; Type: SEQUENCE OWNED BY; Schema: datacatalog; Owner: postgres_user
--

ALTER SEQUENCE datacatalog.schema_id_seq OWNED BY datacatalog.schema_catalogue.id;


--
-- Name: table_catalogue; Type: TABLE; Schema: datacatalog; Owner: postgres_user
--

CREATE TABLE datacatalog.table_catalogue (
    id bigint NOT NULL,
    name character varying NOT NULL,
    schema_id bigint NOT NULL,
    table_type_id bigint NOT NULL,
    description character varying NOT NULL,
    domain_id bigint NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at time with time zone DEFAULT now() NOT NULL,
    updated_at time with time zone DEFAULT now() NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE datacatalog.table_catalogue OWNER TO postgres_user;

--
-- Name: COLUMN table_catalogue.user_id; Type: COMMENT; Schema: datacatalog; Owner: postgres_user
--

COMMENT ON COLUMN datacatalog.table_catalogue.user_id IS 'Updated by';


--
-- Name: table_properties_id_seq; Type: SEQUENCE; Schema: datacatalog; Owner: postgres_user
--

CREATE SEQUENCE datacatalog.table_properties_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE datacatalog.table_properties_id_seq OWNER TO postgres_user;

--
-- Name: table_properties_id_seq; Type: SEQUENCE OWNED BY; Schema: datacatalog; Owner: postgres_user
--

ALTER SEQUENCE datacatalog.table_properties_id_seq OWNED BY datacatalog.table_catalogue.id;


--
-- Name: table_type_catalogue; Type: TABLE; Schema: datacatalog; Owner: postgres_user
--

CREATE TABLE datacatalog.table_type_catalogue (
    id bigint NOT NULL,
    name character varying NOT NULL,
    description character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at time with time zone DEFAULT now() NOT NULL,
    updated_at time with time zone DEFAULT now() NOT NULL,
    updated_by bigint NOT NULL
);


ALTER TABLE datacatalog.table_type_catalogue OWNER TO postgres_user;

--
-- Name: COLUMN table_type_catalogue.updated_by; Type: COMMENT; Schema: datacatalog; Owner: postgres_user
--

COMMENT ON COLUMN datacatalog.table_type_catalogue.updated_by IS 'User Id';


--
-- Name: table_type_id_seq; Type: SEQUENCE; Schema: datacatalog; Owner: postgres_user
--

CREATE SEQUENCE datacatalog.table_type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE datacatalog.table_type_id_seq OWNER TO postgres_user;

--
-- Name: table_type_id_seq; Type: SEQUENCE OWNED BY; Schema: datacatalog; Owner: postgres_user
--

ALTER SEQUENCE datacatalog.table_type_id_seq OWNED BY datacatalog.table_type_catalogue.id;


--
-- Name: action; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.action (
    id bigint NOT NULL,
    name character varying NOT NULL,
    description character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE users.action OWNER TO postgres_user;

--
-- Name: TABLE action; Type: COMMENT; Schema: users; Owner: postgres_user
--

COMMENT ON TABLE users.action IS 'Actions that user or role can do';


--
-- Name: actions_access_id_seq; Type: SEQUENCE; Schema: users; Owner: postgres_user
--

CREATE SEQUENCE users.actions_access_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE users.actions_access_id_seq OWNER TO postgres_user;

--
-- Name: actions_access_id_seq; Type: SEQUENCE OWNED BY; Schema: users; Owner: postgres_user
--

ALTER SEQUENCE users.actions_access_id_seq OWNED BY users.action.id;


--
-- Name: mail_token; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.mail_token (
    hash bytea NOT NULL,
    user_id bigint NOT NULL,
    expire timestamp with time zone NOT NULL,
    scope text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone,
    is_active boolean DEFAULT true NOT NULL
);


ALTER TABLE users.mail_token OWNER TO postgres_user;

--
-- Name: refresh_token; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.refresh_token (
    user_id bigint NOT NULL,
    expire timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    refresh_token character varying(512) NOT NULL,
    id character varying(256) NOT NULL
);


ALTER TABLE users.refresh_token OWNER TO postgres_user;

--
-- Name: resource; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.resource (
    id bigint NOT NULL,
    name character varying NOT NULL,
    description character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE users.resource OWNER TO postgres_user;

--
-- Name: TABLE resource; Type: COMMENT; Schema: users; Owner: postgres_user
--

COMMENT ON TABLE users.resource IS 'Type of resource to be accessed';


--
-- Name: resource_id_seq; Type: SEQUENCE; Schema: users; Owner: postgres_user
--

CREATE SEQUENCE users.resource_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE users.resource_id_seq OWNER TO postgres_user;

--
-- Name: resource_id_seq; Type: SEQUENCE OWNED BY; Schema: users; Owner: postgres_user
--

ALTER SEQUENCE users.resource_id_seq OWNED BY users.resource.id;


--
-- Name: resource_type; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.resource_type (
    id bigint NOT NULL,
    name character varying NOT NULL,
    description character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE users.resource_type OWNER TO postgres_user;

--
-- Name: resource_type_id_seq; Type: SEQUENCE; Schema: users; Owner: postgres_user
--

CREATE SEQUENCE users.resource_type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE users.resource_type_id_seq OWNER TO postgres_user;

--
-- Name: resource_type_id_seq; Type: SEQUENCE OWNED BY; Schema: users; Owner: postgres_user
--

ALTER SEQUENCE users.resource_type_id_seq OWNED BY users.resource_type.id;


--
-- Name: role; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.role (
    id bigint NOT NULL,
    role_name_long character varying NOT NULL,
    role_name_short character varying NOT NULL,
    description character varying,
    jwt_export boolean DEFAULT false NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE users.role OWNER TO postgres_user;

--
-- Name: COLUMN role.jwt_export; Type: COMMENT; Schema: users; Owner: postgres_user
--

COMMENT ON COLUMN users.role.jwt_export IS 'Is it going to be exported to jwt';


--
-- Name: role_access; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.role_access (
    id bigint NOT NULL,
    role_id bigint NOT NULL,
    resource_id bigint NOT NULL,
    action_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    resource_type_id bigint NOT NULL
);


ALTER TABLE users.role_access OWNER TO postgres_user;

--
-- Name: roles_access_id_seq; Type: SEQUENCE; Schema: users; Owner: postgres_user
--

CREATE SEQUENCE users.roles_access_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE users.roles_access_id_seq OWNER TO postgres_user;

--
-- Name: roles_access_id_seq; Type: SEQUENCE OWNED BY; Schema: users; Owner: postgres_user
--

ALTER SEQUENCE users.roles_access_id_seq OWNED BY users.role_access.id;


--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: users; Owner: postgres_user
--

CREATE SEQUENCE users.roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE users.roles_id_seq OWNER TO postgres_user;

--
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: users; Owner: postgres_user
--

ALTER SEQUENCE users.roles_id_seq OWNED BY users.role.id;


--
-- Name: user; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users."user" (
    id bigint NOT NULL,
    name character varying(512) NOT NULL,
    email character varying(512) NOT NULL,
    password_hash bytea NOT NULL,
    activated boolean,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE users."user" OWNER TO postgres_user;

--
-- Name: user_access; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.user_access (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    resource_id bigint NOT NULL,
    action_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    resource_type_id bigint NOT NULL
);


ALTER TABLE users.user_access OWNER TO postgres_user;

--
-- Name: user_access_id_seq; Type: SEQUENCE; Schema: users; Owner: postgres_user
--

CREATE SEQUENCE users.user_access_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE users.user_access_id_seq OWNER TO postgres_user;

--
-- Name: user_access_id_seq; Type: SEQUENCE OWNED BY; Schema: users; Owner: postgres_user
--

ALTER SEQUENCE users.user_access_id_seq OWNED BY users.user_access.id;


--
-- Name: user_role; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.user_role (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    role_id bigint NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE users.user_role OWNER TO postgres_user;

--
-- Name: user_role_id_seq; Type: SEQUENCE; Schema: users; Owner: postgres_user
--

CREATE SEQUENCE users.user_role_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE users.user_role_id_seq OWNER TO postgres_user;

--
-- Name: user_role_id_seq; Type: SEQUENCE OWNED BY; Schema: users; Owner: postgres_user
--

ALTER SEQUENCE users.user_role_id_seq OWNED BY users.user_role.id;


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: users; Owner: postgres_user
--

CREATE SEQUENCE users.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE users.users_id_seq OWNER TO postgres_user;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: users; Owner: postgres_user
--

ALTER SEQUENCE users.users_id_seq OWNED BY users."user".id;


--
-- Name: alias_catalogue id; Type: DEFAULT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.alias_catalogue ALTER COLUMN id SET DEFAULT nextval('datacatalog.alias_id_seq'::regclass);


--
-- Name: column_catalogue id; Type: DEFAULT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.column_catalogue ALTER COLUMN id SET DEFAULT nextval('datacatalog.column_catalogue_id_seq'::regclass);


--
-- Name: database_catalogue id; Type: DEFAULT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.database_catalogue ALTER COLUMN id SET DEFAULT nextval('datacatalog.database_id_seq'::regclass);


--
-- Name: database_type_catalogue id; Type: DEFAULT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.database_type_catalogue ALTER COLUMN id SET DEFAULT nextval('datacatalog.database_type_id_seq'::regclass);


--
-- Name: domain_catalogue id; Type: DEFAULT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.domain_catalogue ALTER COLUMN id SET DEFAULT nextval('datacatalog.domain_id_seq'::regclass);


--
-- Name: group_levels_catalogue id; Type: DEFAULT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.group_levels_catalogue ALTER COLUMN id SET DEFAULT nextval('datacatalog.group_levels_catalogue_id_seq'::regclass);


--
-- Name: host_catalogue id; Type: DEFAULT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.host_catalogue ALTER COLUMN id SET DEFAULT nextval('datacatalog.host_id_seq'::regclass);


--
-- Name: schema_catalogue id; Type: DEFAULT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.schema_catalogue ALTER COLUMN id SET DEFAULT nextval('datacatalog.schema_id_seq'::regclass);


--
-- Name: table_catalogue id; Type: DEFAULT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.table_catalogue ALTER COLUMN id SET DEFAULT nextval('datacatalog.table_properties_id_seq'::regclass);


--
-- Name: table_type_catalogue id; Type: DEFAULT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.table_type_catalogue ALTER COLUMN id SET DEFAULT nextval('datacatalog.table_type_id_seq'::regclass);


--
-- Name: action id; Type: DEFAULT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.action ALTER COLUMN id SET DEFAULT nextval('users.actions_access_id_seq'::regclass);


--
-- Name: resource id; Type: DEFAULT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.resource ALTER COLUMN id SET DEFAULT nextval('users.resource_id_seq'::regclass);


--
-- Name: resource_type id; Type: DEFAULT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.resource_type ALTER COLUMN id SET DEFAULT nextval('users.resource_type_id_seq'::regclass);


--
-- Name: role id; Type: DEFAULT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.role ALTER COLUMN id SET DEFAULT nextval('users.roles_id_seq'::regclass);


--
-- Name: role_access id; Type: DEFAULT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.role_access ALTER COLUMN id SET DEFAULT nextval('users.roles_access_id_seq'::regclass);


--
-- Name: user id; Type: DEFAULT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users."user" ALTER COLUMN id SET DEFAULT nextval('users.users_id_seq'::regclass);


--
-- Name: user_access id; Type: DEFAULT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.user_access ALTER COLUMN id SET DEFAULT nextval('users.user_access_id_seq'::regclass);


--
-- Name: user_role id; Type: DEFAULT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.user_role ALTER COLUMN id SET DEFAULT nextval('users.user_role_id_seq'::regclass);


--
-- Name: alias_catalogue alias_catalogue_unique; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.alias_catalogue
    ADD CONSTRAINT alias_catalogue_unique UNIQUE (name);


--
-- Name: alias_catalogue alias_pk; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.alias_catalogue
    ADD CONSTRAINT alias_pk PRIMARY KEY (id);


--
-- Name: column_catalogue column_catalogue_pk; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.column_catalogue
    ADD CONSTRAINT column_catalogue_pk PRIMARY KEY (id);


--
-- Name: database_catalogue database_pk; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.database_catalogue
    ADD CONSTRAINT database_pk PRIMARY KEY (id);


--
-- Name: database_type_catalogue database_type_pk; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.database_type_catalogue
    ADD CONSTRAINT database_type_pk PRIMARY KEY (id);


--
-- Name: domain_catalogue domain_name_unique; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.domain_catalogue
    ADD CONSTRAINT domain_name_unique UNIQUE (name);


--
-- Name: domain_catalogue domain_pk; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.domain_catalogue
    ADD CONSTRAINT domain_pk PRIMARY KEY (id);


--
-- Name: group_levels_catalogue group_levels_catalogue_pk; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.group_levels_catalogue
    ADD CONSTRAINT group_levels_catalogue_pk PRIMARY KEY (id);


--
-- Name: group_levels_catalogue group_levels_catalogue_unique; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.group_levels_catalogue
    ADD CONSTRAINT group_levels_catalogue_unique UNIQUE (parent_column_id, column_id);


--
-- Name: host_catalogue host_pk; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.host_catalogue
    ADD CONSTRAINT host_pk PRIMARY KEY (id);


--
-- Name: schema_catalogue schema_pk; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.schema_catalogue
    ADD CONSTRAINT schema_pk PRIMARY KEY (id);


--
-- Name: table_catalogue table_properties_pk; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.table_catalogue
    ADD CONSTRAINT table_properties_pk PRIMARY KEY (id);


--
-- Name: table_type_catalogue table_type_name_unique; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.table_type_catalogue
    ADD CONSTRAINT table_type_name_unique UNIQUE (name);


--
-- Name: table_type_catalogue table_type_pk; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.table_type_catalogue
    ADD CONSTRAINT table_type_pk PRIMARY KEY (id);


--
-- Name: action action_access_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.action
    ADD CONSTRAINT action_access_pk PRIMARY KEY (id);


--
-- Name: action action_access_unique; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.action
    ADD CONSTRAINT action_access_unique UNIQUE (name);


--
-- Name: refresh_token refresh_token_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.refresh_token
    ADD CONSTRAINT refresh_token_pk PRIMARY KEY (id);


--
-- Name: resource resource_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.resource
    ADD CONSTRAINT resource_pk PRIMARY KEY (id);


--
-- Name: resource_type resource_type_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.resource_type
    ADD CONSTRAINT resource_type_pk PRIMARY KEY (id);


--
-- Name: resource_type resource_type_unique; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.resource_type
    ADD CONSTRAINT resource_type_unique UNIQUE (name);


--
-- Name: resource resource_unique; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.resource
    ADD CONSTRAINT resource_unique UNIQUE (name);


--
-- Name: role role_name_long_unique; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.role
    ADD CONSTRAINT role_name_long_unique UNIQUE (role_name_long);


--
-- Name: role role_name_short_unique; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.role
    ADD CONSTRAINT role_name_short_unique UNIQUE (role_name_short);


--
-- Name: role role_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.role
    ADD CONSTRAINT role_pk PRIMARY KEY (id);


--
-- Name: role_access roles_access_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.role_access
    ADD CONSTRAINT roles_access_pk PRIMARY KEY (id);


--
-- Name: mail_token token_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.mail_token
    ADD CONSTRAINT token_pk PRIMARY KEY (hash);


--
-- Name: user_access user_access_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.user_access
    ADD CONSTRAINT user_access_pk PRIMARY KEY (id);


--
-- Name: user user_email_unique; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users."user"
    ADD CONSTRAINT user_email_unique UNIQUE (email);


--
-- Name: user user_name_unique; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users."user"
    ADD CONSTRAINT user_name_unique UNIQUE (name);


--
-- Name: user user_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users."user"
    ADD CONSTRAINT user_pk PRIMARY KEY (id);


--
-- Name: user_role user_role_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.user_role
    ADD CONSTRAINT user_role_pk PRIMARY KEY (id);


--
-- Name: column_catalogue_alias_id_idx; Type: INDEX; Schema: datacatalog; Owner: postgres_user
--

CREATE INDEX column_catalogue_alias_id_idx ON datacatalog.column_catalogue USING btree (alias_id);


--
-- Name: column_catalogue_name_idx; Type: INDEX; Schema: datacatalog; Owner: postgres_user
--

CREATE UNIQUE INDEX column_catalogue_name_idx ON datacatalog.column_catalogue USING btree (name);


--
-- Name: table_schema_id_idx; Type: INDEX; Schema: datacatalog; Owner: postgres_user
--

CREATE UNIQUE INDEX table_schema_id_idx ON datacatalog.table_catalogue USING btree (schema_id, name);


--
-- Name: role_access_resource_type_id_idx; Type: INDEX; Schema: users; Owner: postgres_user
--

CREATE INDEX role_access_resource_type_id_idx ON users.role_access USING btree (resource_type_id, resource_id, role_id);


--
-- Name: role_access_role_id_idx; Type: INDEX; Schema: users; Owner: postgres_user
--

CREATE UNIQUE INDEX role_access_role_id_idx ON users.role_access USING btree (role_id, resource_type_id, resource_id);


--
-- Name: user_access_resource_type_id_idx; Type: INDEX; Schema: users; Owner: postgres_user
--

CREATE INDEX user_access_resource_type_id_idx ON users.user_access USING btree (resource_type_id, resource_id, user_id);


--
-- Name: user_access_user_id_idx; Type: INDEX; Schema: users; Owner: postgres_user
--

CREATE INDEX user_access_user_id_idx ON users.user_access USING btree (user_id, resource_type_id, resource_id);


--
-- Name: alias_catalogue alias_user_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.alias_catalogue
    ADD CONSTRAINT alias_user_fk FOREIGN KEY (user_id) REFERENCES users."user"(id);


--
-- Name: column_catalogue column_catalogue_alias_catalogue_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.column_catalogue
    ADD CONSTRAINT column_catalogue_alias_catalogue_fk FOREIGN KEY (alias_id) REFERENCES datacatalog.alias_catalogue(id);


--
-- Name: column_catalogue column_catalogue_table_catalogue_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.column_catalogue
    ADD CONSTRAINT column_catalogue_table_catalogue_fk FOREIGN KEY (table_id) REFERENCES datacatalog.table_catalogue(id);


--
-- Name: column_catalogue column_catalogue_user_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.column_catalogue
    ADD CONSTRAINT column_catalogue_user_fk FOREIGN KEY (user_id) REFERENCES users."user"(id);


--
-- Name: database_catalogue database_host_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.database_catalogue
    ADD CONSTRAINT database_host_fk FOREIGN KEY (host_id) REFERENCES datacatalog.host_catalogue(id);


--
-- Name: database_type_catalogue database_type_user_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.database_type_catalogue
    ADD CONSTRAINT database_type_user_fk FOREIGN KEY (user_id) REFERENCES users."user"(id);


--
-- Name: database_catalogue database_user_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.database_catalogue
    ADD CONSTRAINT database_user_fk FOREIGN KEY (user_id) REFERENCES users."user"(id);


--
-- Name: group_levels_catalogue group_levels_catalogue_column_catalogue_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.group_levels_catalogue
    ADD CONSTRAINT group_levels_catalogue_column_catalogue_fk FOREIGN KEY (column_id) REFERENCES datacatalog.column_catalogue(id);


--
-- Name: group_levels_catalogue group_levels_catalogue_column_catalogue_fk_1; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.group_levels_catalogue
    ADD CONSTRAINT group_levels_catalogue_column_catalogue_fk_1 FOREIGN KEY (parent_column_id) REFERENCES datacatalog.column_catalogue(id);


--
-- Name: group_levels_catalogue group_levels_catalogue_user_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.group_levels_catalogue
    ADD CONSTRAINT group_levels_catalogue_user_fk FOREIGN KEY (user_id) REFERENCES users."user"(id);


--
-- Name: host_catalogue host_database_type_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.host_catalogue
    ADD CONSTRAINT host_database_type_fk FOREIGN KEY (database_type_id) REFERENCES datacatalog.database_type_catalogue(id);


--
-- Name: host_catalogue host_user_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.host_catalogue
    ADD CONSTRAINT host_user_fk FOREIGN KEY (user_id) REFERENCES users."user"(id);


--
-- Name: schema_catalogue schema_database_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.schema_catalogue
    ADD CONSTRAINT schema_database_fk FOREIGN KEY (database_id) REFERENCES datacatalog.database_catalogue(id);


--
-- Name: schema_catalogue schema_user_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.schema_catalogue
    ADD CONSTRAINT schema_user_fk FOREIGN KEY (user_id) REFERENCES users."user"(id);


--
-- Name: table_catalogue table_properties_domain_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.table_catalogue
    ADD CONSTRAINT table_properties_domain_fk FOREIGN KEY (domain_id) REFERENCES datacatalog.domain_catalogue(id);


--
-- Name: table_catalogue table_properties_table_type_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.table_catalogue
    ADD CONSTRAINT table_properties_table_type_fk FOREIGN KEY (table_type_id) REFERENCES datacatalog.table_type_catalogue(id);


--
-- Name: table_catalogue table_properties_user_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.table_catalogue
    ADD CONSTRAINT table_properties_user_fk FOREIGN KEY (user_id) REFERENCES users."user"(id);


--
-- Name: table_catalogue table_schema_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.table_catalogue
    ADD CONSTRAINT table_schema_fk FOREIGN KEY (schema_id) REFERENCES datacatalog.schema_catalogue(id);


--
-- Name: table_type_catalogue table_type_user_fk; Type: FK CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.table_type_catalogue
    ADD CONSTRAINT table_type_user_fk FOREIGN KEY (updated_by) REFERENCES users."user"(id);


--
-- Name: mail_token mail_token_users_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.mail_token
    ADD CONSTRAINT mail_token_users_fk FOREIGN KEY (user_id) REFERENCES users."user"(id) ON DELETE CASCADE;


--
-- Name: refresh_token refresh_token_users_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.refresh_token
    ADD CONSTRAINT refresh_token_users_fk FOREIGN KEY (user_id) REFERENCES users."user"(id);


--
-- Name: role_access role_access_action_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.role_access
    ADD CONSTRAINT role_access_action_fk FOREIGN KEY (action_id) REFERENCES users.action(id);


--
-- Name: role_access role_access_resource_type_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.role_access
    ADD CONSTRAINT role_access_resource_type_fk FOREIGN KEY (resource_type_id) REFERENCES users.resource_type(id);


--
-- Name: role_access role_access_role_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.role_access
    ADD CONSTRAINT role_access_role_fk FOREIGN KEY (role_id) REFERENCES users.role(id);


--
-- Name: user_access user_access_users_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.user_access
    ADD CONSTRAINT user_access_users_fk FOREIGN KEY (user_id) REFERENCES users."user"(id);


--
-- Name: user_role user_role_role_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.user_role
    ADD CONSTRAINT user_role_role_fk FOREIGN KEY (role_id) REFERENCES users.role(id);


--
-- Name: user_role user_role_user_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.user_role
    ADD CONSTRAINT user_role_user_fk FOREIGN KEY (user_id) REFERENCES users."user"(id);


--
-- Name: SCHEMA users; Type: ACL; Schema: -; Owner: postgres_user
--

GRANT USAGE ON SCHEMA users TO datacomrade_user;


--
-- Name: TABLE domain_catalogue; Type: ACL; Schema: datacatalog; Owner: postgres_user
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE datacatalog.domain_catalogue TO datacomrade_user;


--
-- Name: TABLE table_catalogue; Type: ACL; Schema: datacatalog; Owner: postgres_user
--

GRANT SELECT,INSERT,UPDATE ON TABLE datacatalog.table_catalogue TO datacomrade_user;


--
-- Name: TABLE table_type_catalogue; Type: ACL; Schema: datacatalog; Owner: postgres_user
--

GRANT SELECT,INSERT,UPDATE ON TABLE datacatalog.table_type_catalogue TO datacomrade_user;


--
-- Name: TABLE action; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT SELECT,INSERT,UPDATE ON TABLE users.action TO datacomrade_user;


--
-- Name: TABLE mail_token; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT SELECT,INSERT,REFERENCES,DELETE,TRIGGER,MAINTAIN,UPDATE ON TABLE users.mail_token TO datacomrade_user;


--
-- Name: TABLE refresh_token; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT ALL ON TABLE users.refresh_token TO datacomrade_user;


--
-- Name: TABLE resource; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT SELECT,INSERT,UPDATE ON TABLE users.resource TO datacomrade_user;


--
-- Name: TABLE resource_type; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT SELECT,INSERT,UPDATE ON TABLE users.resource_type TO datacomrade_user;


--
-- Name: TABLE role; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT SELECT,INSERT,REFERENCES,DELETE,TRIGGER,MAINTAIN,UPDATE ON TABLE users.role TO datacomrade_user;


--
-- Name: TABLE role_access; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT SELECT,INSERT,UPDATE ON TABLE users.role_access TO datacomrade_user;


--
-- Name: SEQUENCE roles_id_seq; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT SELECT,USAGE ON SEQUENCE users.roles_id_seq TO datacomrade_user;


--
-- Name: TABLE "user"; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT ALL ON TABLE users."user" TO datacomrade_user;


--
-- Name: COLUMN "user".name; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT INSERT(name) ON TABLE users."user" TO datacomrade_user;


--
-- Name: TABLE user_access; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT SELECT,INSERT,UPDATE ON TABLE users.user_access TO datacomrade_user;


--
-- Name: TABLE user_role; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT SELECT,INSERT,UPDATE ON TABLE users.user_role TO datacomrade_user;


--
-- Name: SEQUENCE users_id_seq; Type: ACL; Schema: users; Owner: postgres_user
--

GRANT SELECT,USAGE ON SEQUENCE users.users_id_seq TO datacomrade_user;


--
-- PostgreSQL database dump complete
--

