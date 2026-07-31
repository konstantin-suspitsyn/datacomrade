-- +goose up

-- +goose StatementBegin
CREATE FUNCTION dc.get_table_ids_by_external_user_and_role(
    p_external_id text,
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
-- +goose StatementEnd

-- +goose down

DROP FUNCTION IF EXISTS dc.get_table_ids_by_external_user_and_role(text, text);
