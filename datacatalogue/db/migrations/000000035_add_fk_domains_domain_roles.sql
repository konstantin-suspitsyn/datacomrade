-- +goose up
ALTER TABLE dc.domains_domain_roles ADD CONSTRAINT domains_domain_roles_domain_cat_fk FOREIGN KEY (domain_cat_id) REFERENCES dc.domain_cat(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.domains_domain_roles ADD CONSTRAINT domains_domain_roles_domain_roles_fk FOREIGN KEY (domain_roles_id) REFERENCES dc.domain_roles(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.domains_domain_roles DROP CONSTRAINT domains_domain_roles_domain_cat_fk;
ALTER TABLE dc.domains_domain_roles DROP CONSTRAINT domains_domain_roles_domain_roles_fk;