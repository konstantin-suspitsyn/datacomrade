-- +goose up
ALTER TABLE dc.user_domain_roles ADD CONSTRAINT user_domain_roles_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.user_domain_roles ADD CONSTRAINT user_domain_roles_domain_roles_fk FOREIGN KEY (domain_roles_id) REFERENCES dc.domain_roles(id);

-- +goose down
ALTER TABLE dc.user_domain_roles DROP CONSTRAINT user_domain_roles_user_fk;
ALTER TABLE dc.user_domain_roles DROP CONSTRAINT user_domain_roles_domain_roles_fk;