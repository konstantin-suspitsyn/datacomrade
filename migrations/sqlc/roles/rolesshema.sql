CREATE SCHEMA users; 

CREATE TABLE users.role_access (
	id int8 DEFAULT nextval('users.roles_access_id_seq'::regclass) NOT NULL,
	role_id int8 NOT NULL,
	resource_id int8 NOT NULL,
	action_id int8 NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	resource_type_id int8 NOT NULL,
	CONSTRAINT roles_access_pk PRIMARY KEY (id)
);

CREATE TABLE users.user_access (
	id bigserial NOT NULL,
	user_id int8 NOT NULL,
	resource_id int8 NOT NULL,
	action_id int8 NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	resource_type_id int8 NOT NULL,
	CONSTRAINT user_access_pk PRIMARY KEY (id)
);

CREATE TABLE users.user_role (
	id bigserial NOT NULL,
	user_id int8 NOT NULL,
	role_id int8 NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NULL,
	CONSTRAINT user_role_pk PRIMARY KEY (id)
);

CREATE TABLE users."role" (
	id int8 DEFAULT nextval('users.roles_id_seq'::regclass) NOT NULL,
	role_name_long varchar NOT NULL,
	role_name_short varchar NOT NULL,
	description varchar NULL,
	jwt_export bool DEFAULT false NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT role_name_long_unique UNIQUE (role_name_long),
	CONSTRAINT role_name_short_unique UNIQUE (role_name_short),
	CONSTRAINT role_pk PRIMARY KEY (id)
);

CREATE TABLE users.resource (
	id bigserial NOT NULL,
	"name" varchar NOT NULL,
	description varchar NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT resource_pk PRIMARY KEY (id),
	CONSTRAINT resource_unique UNIQUE (name)
);

CREATE TABLE users."action" (
	id int8 DEFAULT nextval('users.actions_access_id_seq'::regclass) NOT NULL,
	"name" varchar NOT NULL,
	description varchar NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT action_access_pk PRIMARY KEY (id),
	CONSTRAINT action_access_unique UNIQUE (name)
);
