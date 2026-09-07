CREATE SCHEMA dc;

-- dc.user не участвует в самом запросе напрямую — она видна только внутри
-- dc.get_table_ids_by_external_user_and_role, но её колонка external_id
-- нужна конвейеру, чтобы найти границу и тип для одноимённого параметра
-- функции (см. auth_logic/schema.sql — та же таблица, та же причина).
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

CREATE TABLE dc.alias (
                          id bigserial NOT NULL,
                          "name" varchar(255) NOT NULL,
                          description varchar(1000) NOT NULL,
                          created_at timestamp DEFAULT now() NOT NULL,
                          updated_at timestamp DEFAULT now() NULL,
                          is_deleted bool DEFAULT false NOT NULL,
                          user_id int8 NOT NULL,
                          column_type_id int8 DEFAULT 1 NOT NULL
);

CREATE TABLE dc.column_type (
                                id bigserial NOT NULL,
                                "name" varchar(128) NOT NULL,
                                description varchar(1000) NOT NULL,
                                is_deleted bool DEFAULT false NOT NULL,
                                created_at timestamp DEFAULT now() NOT NULL,
                                updated_at timestamp DEFAULT now() NOT NULL,
                                user_id int8 NOT NULL
);
CREATE TABLE dc.column_cat (
                               id bigserial NOT NULL,
                               table_id int8 NOT NULL,
                               "name" varchar(256) NOT NULL,
                               alias_id int8 NULL,
                               column_type_id int8 NOT NULL,
                               description varchar(1000) NOT NULL,
                               calculation_type_id int8 NULL,
                               is_deleted bool DEFAULT false NOT NULL,
                               show_in_ui bool NOT NULL,
                               created_at timestamp DEFAULT now() NOT NULL,
                               updated_at timestamp DEFAULT now() NOT NULL,
                               user_id int8 NOT NULL
);

CREATE FUNCTION dc.get_table_ids_by_external_user_and_role(
    p_external_id uuid,
    p_role        text
)
    RETURNS TABLE (id bigint)
    LANGUAGE plpgsql
STABLE
AS $$
BEGIN
RETURN QUERY
    WITH filtered_user AS (
        SELECT u.id
        FROM dc."user" u
        WHERE u.external_id = p_external_id::uuid
    ),
    all_domains AS (
        -- Все домены, которые пользователь может прочитать
        SELECT udr.domain_id
        FROM dc.user_domain_roles udr
        INNER JOIN dc.domain_roles dr ON dr.id = udr.domain_roles_id
        WHERE udr.user_id = (SELECT fu.id FROM filtered_user fu)
          AND dr.name = p_role
    ),
    all_tables_with_role_and_user AS (
        -- Все таблицы, которые пользователь может прочитать
        SELECT utr.table_id
        FROM dc.user_table_roles utr
        JOIN dc.table_roles tr ON tr.id = utr.table_roles_id
        WHERE tr.name = p_role
          AND utr.user_id = (SELECT fu.id FROM filtered_user fu)
    )
SELECT tc.id
FROM dc.table_cat tc
WHERE tc.domain_id IN (SELECT domain_id FROM all_domains)
   OR tc.id IN (SELECT table_id FROM all_tables_with_role_and_user);
END;
$$;

CREATE TABLE dc.table_cat (
                              id bigserial NOT NULL,
                              "name" varchar(128) NOT NULL,
                              description varchar(2000) NOT NULL,
                              schema_id int8 NOT NULL,
                              table_type_id int8 NOT NULL,
                              domain_id int8 NOT NULL,
                              is_deleted bool DEFAULT false NOT NULL,
                              is_get_dict bool DEFAULT false NOT NULL,
                              created_at timestamp DEFAULT now() NOT NULL,
                              updated_at timestamp DEFAULT now() NOT NULL,
                              user_id int8 NOT NULL
);
CREATE TABLE dc.domain_cat (
                               id bigserial NOT NULL,
                               domain_name varchar(100) NOT NULL,
                               is_deleted bool DEFAULT false NOT NULL,
                               created_at timestamp DEFAULT now() NOT NULL,
                               updated_at timestamp DEFAULT now() NOT NULL,
                               user_id int8 NOT NULL
);
CREATE TABLE dc.table_type (
                               id bigserial NOT NULL,
                               "name" varchar(128) NOT NULL,
                               description varchar(1000) NOT NULL,
                               is_deleted bool DEFAULT false NOT NULL,
                               created_at timestamp DEFAULT now() NOT NULL,
                               updated_at timestamp DEFAULT now() NOT NULL,
                               user_id int8 NOT NULL
);
CREATE TABLE dc.schema_cat (
                               id bigserial NOT NULL,
                               database_id int8 NOT NULL,
                               "name" varchar(128) NOT NULL,
                               is_deleted bool DEFAULT false NOT NULL,
                               created_at timestamp DEFAULT now() NOT NULL,
                               updated_at timestamp DEFAULT now() NOT NULL,
                               user_id int8 NOT NULL
);
CREATE TABLE dc.database_cat (
                                 id bigserial NOT NULL,
                                 "name" varchar(255) NOT NULL,
                                 host_id int8 NOT NULL,
                                 database_type_id int8 NOT NULL,
                                 description varchar(1000) NOT NULL,
                                 is_deleted bool DEFAULT false NOT NULL,
                                 created_at timestamp DEFAULT now() NOT NULL,
                                 updated_at timestamp DEFAULT now() NOT NULL,
                                 user_id int8 NOT NULL
);
CREATE TABLE dc.database_type (
                                  id bigserial NOT NULL,
                                  "name" varchar(128) NOT NULL,
                                  db_version varchar(512) NOT NULL,
                                  is_deleted bool DEFAULT false NOT NULL,
                                  created_at timestamp DEFAULT now() NOT NULL,
                                  updated_at timestamp DEFAULT now() NOT NULL,
                                  user_id int8 NOT NULL
);