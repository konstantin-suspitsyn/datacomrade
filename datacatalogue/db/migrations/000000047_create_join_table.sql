-- +goose up
CREATE TABLE dc.tables_joins (
                                 id bigserial NOT NULL,
                                 table_cat_id_left bigint NOT NULL,
                                 table_cat_id_right bigint NOT NULL,
                                 is_deleted boolean DEFAULT false NOT NULL,
                                 created_at timestamp DEFAULT now() NOT NULL,
                                 updated_at timestamp DEFAULT now() NOT NULL,
                                 user_id bigint NOT NULL,
                                 CONSTRAINT tables_joins_pk PRIMARY KEY (id),
                                 CONSTRAINT tables_joins_unique UNIQUE (table_cat_id_left,table_cat_id_right),
                                 CONSTRAINT table_cat_left_fk FOREIGN KEY (table_cat_id_left) REFERENCES dc.table_cat(id),
                                 CONSTRAINT table_cat_right_fk FOREIGN KEY (table_cat_id_right) REFERENCES dc.table_cat(id),
                                 CONSTRAINT tables_joins_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

-- +goose down

DROP TABLE dc.tables_joins;
