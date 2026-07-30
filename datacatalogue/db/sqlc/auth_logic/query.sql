-- =========================================================
-- Логика авторизации. А именно доступные поля и таблицы
-- =========================================================

-- =========================================================
-- Таблицы, доступные пользователю
-- =========================================================

-- name: GetTableIdsByExternalUserIdAndRoles :many
with filtered_user as (select id
                       from dc."user"
                       where external_id = $1),
     all_domains as (
         -- Все домены, которые пользователь может прочитать
         select udr.domain_id
         from dc.user_domain_roles udr
                  inner join dc.domain_roles dr on dr.id = udr.domain_roles_id
         where udr.user_id = (select id from filtered_user)
           and dr.name in ($2))
        ,
     all_tables_with_role_and_user as (
         -- Все таблицы, которые пользователь может прочитать
         select utr.table_id
         from dc.user_table_roles utr
                  join dc.table_roles tr on tr.id = utr.table_roles_id
         where tr.name in ($2)
           and utr.user_id = (select id from filtered_user))

select tc.id
from dc.table_cat tc
where tc.domain_id in (select domain_id from all_domains)
   or tc.id in (select table_id from all_tables_with_role_and_user);

-- name: GetTableIdsByUserIdAndRoles :many
with all_domains as (
         -- Все домены, которые пользователь может прочитать
         select udr.domain_id
         from dc.user_domain_roles udr
                  inner join dc.domain_roles dr on dr.id = udr.domain_roles_id
         where udr.user_id = $1
           and dr.name in ($2))
        ,
     all_tables_with_role_and_user as (
         -- Все таблицы, которые пользователь может прочитать
         select utr.table_id
         from dc.user_table_roles utr
                  join dc.table_roles tr on tr.id = utr.table_roles_id
         where tr.name in ($2)
           and utr.user_id = $1)

select tc.id
from dc.table_cat tc
where tc.domain_id in (select domain_id from all_domains)
   or tc.id in (select table_id from all_tables_with_role_and_user);

