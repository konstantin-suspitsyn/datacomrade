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
-- Name: domain; Type: TABLE; Schema: datacatalog; Owner: postgres_user
--

CREATE TABLE datacatalog.domain (
    id bigint NOT NULL,
    name character varying NOT NULL,
    description character varying
);


ALTER TABLE datacatalog.domain OWNER TO postgres_user;

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

ALTER SEQUENCE datacatalog.domain_id_seq OWNED BY datacatalog.domain.id;


--
-- Name: domain_role; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.domain_role (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    domain_id bigint NOT NULL,
    role_id bigint NOT NULL
);


ALTER TABLE users.domain_role OWNER TO postgres_user;

--
-- Name: domain_role_id_seq; Type: SEQUENCE; Schema: users; Owner: postgres_user
--

CREATE SEQUENCE users.domain_role_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE users.domain_role_id_seq OWNER TO postgres_user;

--
-- Name: domain_role_id_seq; Type: SEQUENCE OWNED BY; Schema: users; Owner: postgres_user
--

ALTER SEQUENCE users.domain_role_id_seq OWNED BY users.domain_role.id;


--
-- Name: mail_token; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.mail_token (
    hash bytea NOT NULL,
    user_id bigint NOT NULL,
    expire timestamp with time zone NOT NULL,
    scope text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
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
    updated_at timestamp without time zone DEFAULT now(),
    refresh_token character varying(512) NOT NULL,
    id character varying(256) NOT NULL
);


ALTER TABLE users.refresh_token OWNER TO postgres_user;

--
-- Name: roles; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.roles (
    id bigint NOT NULL,
    role_name_long character varying NOT NULL,
    role_name_short character varying NOT NULL,
    description character varying
);


ALTER TABLE users.roles OWNER TO postgres_user;

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

ALTER SEQUENCE users.roles_id_seq OWNED BY users.roles.id;


--
-- Name: users; Type: TABLE; Schema: users; Owner: postgres_user
--

CREATE TABLE users.users (
    id bigint NOT NULL,
    name character varying(512) NOT NULL,
    email character varying(512) NOT NULL,
    password_hash bytea NOT NULL,
    activated boolean,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE users.users OWNER TO postgres_user;

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

ALTER SEQUENCE users.users_id_seq OWNED BY users.users.id;


--
-- Name: domain id; Type: DEFAULT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.domain ALTER COLUMN id SET DEFAULT nextval('datacatalog.domain_id_seq'::regclass);


--
-- Name: domain_role id; Type: DEFAULT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.domain_role ALTER COLUMN id SET DEFAULT nextval('users.domain_role_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.roles ALTER COLUMN id SET DEFAULT nextval('users.roles_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.users ALTER COLUMN id SET DEFAULT nextval('users.users_id_seq'::regclass);


--
-- Name: domain domain_name_unique; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.domain
    ADD CONSTRAINT domain_name_unique UNIQUE (name);


--
-- Name: domain domain_pk; Type: CONSTRAINT; Schema: datacatalog; Owner: postgres_user
--

ALTER TABLE ONLY datacatalog.domain
    ADD CONSTRAINT domain_pk PRIMARY KEY (id);


--
-- Name: domain_role domain_role_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.domain_role
    ADD CONSTRAINT domain_role_pk PRIMARY KEY (id);


--
-- Name: refresh_token refresh_token_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.refresh_token
    ADD CONSTRAINT refresh_token_pk PRIMARY KEY (id);


--
-- Name: roles roles_name_long_unique; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.roles
    ADD CONSTRAINT roles_name_long_unique UNIQUE (role_name_long);


--
-- Name: roles roles_name_short_unique; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.roles
    ADD CONSTRAINT roles_name_short_unique UNIQUE (role_name_short);


--
-- Name: roles roles_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.roles
    ADD CONSTRAINT roles_pk PRIMARY KEY (id);


--
-- Name: mail_token token_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.mail_token
    ADD CONSTRAINT token_pk PRIMARY KEY (hash);


--
-- Name: users users_email_unique; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.users
    ADD CONSTRAINT users_email_unique UNIQUE (email);


--
-- Name: users users_pk; Type: CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.users
    ADD CONSTRAINT users_pk PRIMARY KEY (id);


--
-- Name: domain_role domain_role_domain_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.domain_role
    ADD CONSTRAINT domain_role_domain_fk FOREIGN KEY (domain_id) REFERENCES datacatalog.domain(id);


--
-- Name: domain_role domain_role_roles_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.domain_role
    ADD CONSTRAINT domain_role_roles_fk FOREIGN KEY (role_id) REFERENCES users.roles(id);


--
-- Name: domain_role domain_role_users_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.domain_role
    ADD CONSTRAINT domain_role_users_fk FOREIGN KEY (user_id) REFERENCES users.users(id);


--
-- Name: mail_token mail_token_users_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.mail_token
    ADD CONSTRAINT mail_token_users_fk FOREIGN KEY (user_id) REFERENCES users.users(id);


--
-- Name: refresh_token refresh_token_users_fk; Type: FK CONSTRAINT; Schema: users; Owner: postgres_user
--

ALTER TABLE ONLY users.refresh_token
    ADD CONSTRAINT refresh_token_users_fk FOREIGN KEY (user_id) REFERENCES users.users(id);


--
-- Name: TABLE domain; Type: ACL; Schema: datacatalog; Owner: postgres_user
--


--
-- Name: TABLE domain_role; Type: ACL; Schema: users; Owner: postgres_user
--



--
-- Name: TABLE roles; Type: ACL; Schema: users; Owner: postgres_user
--



--
-- Name: TABLE users; Type: ACL; Schema: users; Owner: postgres_user
--



--
-- Name: COLUMN users.name; Type: ACL; Schema: users; Owner: postgres_user
--



--
-- PostgreSQL database dump complete
--

