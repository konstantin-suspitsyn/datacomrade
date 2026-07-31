# -*- coding: utf-8 -*-
"""Генерация user_domain_roles.proto из schema.sql/query.sql домена user_domain_roles.

Реализует правила из documentation/dev_instructions/crud/proto_based_on_crud.md:
8 стандартных RPC на таблицу, с спецмаппингом
"<col>_id = (SELECT u.id FROM dc.\"user\" u WHERE u.external_id = $N)"
в proto-поле "<col>_external_id" (см. EXTERNAL_ID_SUFFIX в resolve.py).

Запуск:  python deploy/generators/proto_gen_user_domain_roles.py
"""
import re

SCHEMA_PATH = r"C:\Users\const\working\datacomrade\datacatalogue\db\sqlc\user_domain_roles\schema.sql"
QUERY_PATH = r"C:\Users\const\working\datacomrade\datacatalogue\db\sqlc\user_domain_roles\query.sql"
OUT_PATH = r"C:\Users\const\working\datacomrade\shared\proto\user_domain_roles\v1\user_domain_roles.proto"

ENTITY_NAMES = {
    'domain_roles': 'DomainRole',
    'table_roles': 'TableRole',
    'user_domain_roles': 'UserDomainRole',
    'user_table_roles': 'UserTableRole',
}

TYPE_RE = re.compile(
    r'^(\w+|"[^"]+")\s+('
    r'character varying\(\d+\)|character varying|varchar\(\d+\)'
    r'|timestamp without time zone|timestamptz|timestamp'
    r'|bigint|bigserial|boolean|bool|smallint|integer|int4|int8|int2'
    r'|numeric|uuid|text'
    r')\b', re.IGNORECASE)


def pg_type_to_proto(pg_type):
    t = pg_type.lower()
    if t.startswith('character varying') or t.startswith('varchar') or t == 'text':
        return 'string', False
    if t in ('bigint', 'bigserial'):
        return 'int64', False
    if t in ('boolean', 'bool'):
        return 'bool', False
    if t in ('smallint', 'integer', 'int4', 'int2'):
        return 'int32', False
    if t.startswith('timestamp'):
        return 'google.protobuf.Timestamp', True
    if t == 'uuid':
        return 'string', False
    if t == 'numeric':
        return 'string', False
    raise ValueError(f"unmapped type: {pg_type}")


def split_top_level(s, sep=','):
    parts = []
    depth = 0
    cur = ''
    for ch in s:
        if ch == '(':
            depth += 1
        elif ch == ')':
            depth -= 1
        if ch == sep and depth == 0:
            parts.append(cur.strip())
            cur = ''
        else:
            cur += ch
    if cur.strip():
        parts.append(cur.strip())
    return parts


def parse_schema(path):
    text = open(path, encoding='utf-8').read()
    tables = {}
    for m in re.finditer(r'create table dc\.(\w+)\s*\((.*?)\n\);', text, re.DOTALL | re.IGNORECASE):
        table = m.group(1)
        body = m.group(2)
        cols = []
        for line in split_top_level(body, ','):
            line = line.strip()
            tm = TYPE_RE.match(line)
            if not tm:
                continue
            name = tm.group(1).strip('"')
            pg_type = tm.group(2)
            cols.append((name, pg_type))
        tables[table] = cols
    return tables


def pluralize(word):
    if word.endswith(('s', 'x', 'z')) or word.endswith('ch') or word.endswith('sh'):
        return word + 'es'
    return word + 's'


def entity_snake(entity_pascal):
    return re.sub(r'(?<!^)(?=[A-Z])', '_', entity_pascal).lower()


USER_SUBQUERY_RE = re.compile(
    r'\(SELECT u\.id FROM dc\."user" u WHERE u\.external_id = \$(\d+)\)'
)


def value_to_field(col_name, value, table_cols):
    value = value.strip()
    m = USER_SUBQUERY_RE.match(value)
    if m:
        assert col_name.endswith('_id'), col_name
        return col_name[:-3] + '_external_id', 'string', False
    if re.match(r'^\$\d+$', value):
        pg_type = dict(table_cols)[col_name]
        proto_type, is_msg = pg_type_to_proto(pg_type)
        return col_name, proto_type, is_msg
    return None


def parse_queries(path):
    text = open(path, encoding='utf-8').read()
    section_re = re.compile(
        r'-- =+\n-- (dc\.\w+)\n-- =+\n(.*?)(?=\n-- =+\n-- dc\.|\Z)', re.DOTALL)
    sections = []
    for m in section_re.finditer(text):
        table = m.group(1).split('.', 1)[1]
        body = m.group(2)
        sections.append((table, body))

    query_re = re.compile(r'-- name: (\w+) :(\w+)\n(.*?)(?=\n-- name:|\Z)', re.DOTALL)
    result = []
    for table, body in sections:
        queries = []
        for qm in query_re.finditer(body):
            name, qtype, sql = qm.group(1), qm.group(2), qm.group(3).strip()
            queries.append((name, qtype, sql))
        result.append((table, queries))
    return result


def build_proto():
    schema = parse_schema(SCHEMA_PATH)
    sections = parse_queries(QUERY_PATH)

    service_lines = []
    message_blocks = []

    for idx, (table, queries) in enumerate(sections):
        if idx > 0:
            service_lines.append('')
        entity = ENTITY_NAMES[table]
        e_snake = entity_snake(entity)
        e_plural = pluralize(e_snake)
        entity_plural = pluralize(entity)
        cols = schema[table]

        service_lines.append(f'  // dc.{table}')
        for name, qtype, sql in queries:
            service_lines.append(f'  rpc {name}({name}Request) returns ({name}Response);')

        block = []
        block.append('// =========================================================')
        block.append(f'// dc.{table}')
        block.append('// =========================================================')
        block.append('')
        block.append(f'// {entity} is a full row of dc.{table}.')
        block.append(f'message {entity} {{')
        for i, (cname, ptype) in enumerate(cols, start=1):
            proto_type, _ = pg_type_to_proto(ptype)
            block.append(f'  {proto_type} {cname} = {i};')
        block.append('}')
        block.append('')

        for name, qtype, sql in queries:
            if name == f'Get{entity}ById' or name == f'GetDeleted{entity}ById':
                soft = 'soft deleted' if name.startswith('GetDeleted') else 'active'
                block.append(f'// {name}Request asks for a single {soft} dc.{table} row by id (query {name}).')
                block.append(f'message {name}Request {{')
                block.append('  int64 id = 1;')
                block.append('}')
                block.append('')
                block.append(f'// {name}Response returns a single {soft} dc.{table} row (query {name}).')
                block.append(f'message {name}Response {{')
                block.append(f'  {entity} {e_snake} = 1;')
                block.append('}')
                block.append('')
            elif name == f'Get{entity_plural}' or name == f'GetDeleted{entity_plural}':
                soft = 'soft deleted' if name.startswith('GetDeleted') else 'active'
                block.append(f'// {name}Request asks for every {soft} dc.{table} row (query {name}).')
                block.append(f'message {name}Request {{}}')
                block.append('')
                block.append(f'// {name}Response returns every {soft} dc.{table} row (query {name}).')
                block.append(f'message {name}Response {{')
                block.append(f'  repeated {entity} {e_plural} = 1;')
                block.append('}')
                block.append('')
            elif name == f'Create{entity}':
                m = re.search(r'INSERT INTO dc\.\w+\s*\((.*?)\)\s*VALUES\s*\((.*)\)\s*RETURNING', sql, re.DOTALL)
                col_list = [c.strip().strip('"') for c in split_top_level(m.group(1))]
                val_list = split_top_level(m.group(2))
                fields = []
                for cname, val in zip(col_list, val_list):
                    r = value_to_field(cname, val, cols)
                    if r:
                        fields.append(r)
                block.append(f'// {name}Request carries the fields needed to insert a dc.{table} row (query {name}).')
                block.append(f'message {name}Request {{')
                for i, (fname, ftype, _) in enumerate(fields, start=1):
                    block.append(f'  {ftype} {fname} = {i};')
                block.append('}')
                block.append('')
                block.append(f'// {name}Response returns the created dc.{table} row (query {name}).')
                block.append(f'message {name}Response {{')
                block.append(f'  {entity} {e_snake} = 1;')
                block.append('}')
                block.append('')
            elif name == f'Update{entity}ById':
                m = re.search(r'UPDATE dc\.\w+\s*SET (.*)\nWHERE', sql)
                set_list = split_top_level(m.group(1))
                fields = []
                for assign in set_list:
                    cname, val = assign.split('=', 1)
                    cname = cname.strip().strip('"')
                    r = value_to_field(cname, val, cols)
                    if r:
                        fields.append(r)
                block.append(f'// {name}Request carries the fields to update on a dc.{table} row (query {name}).')
                block.append(f'message {name}Request {{')
                block.append('  int64 id = 1;')
                for i, (fname, ftype, _) in enumerate(fields, start=2):
                    block.append(f'  {ftype} {fname} = {i};')
                block.append('}')
                block.append('')
                block.append(f'// {name}Response returns the updated dc.{table} row (query {name}).')
                block.append(f'message {name}Response {{')
                block.append(f'  {entity} {e_snake} = 1;')
                block.append('}')
                block.append('')
            elif name == f'Delete{entity}ById' or name == f'Undelete{entity}ById':
                if name.startswith('Delete'):
                    action = f'asks to soft delete a dc.{table} row by id, setting is_deleted = true'
                else:
                    action = f'asks to restore a soft deleted dc.{table} row by id, setting is_deleted = false'
                block.append(f'// {name}Request {action} (query {name}).')
                block.append(f'message {name}Request {{')
                block.append('  int64 id = 1;')
                block.append('}')
                block.append('')
                block.append(f'// {name}Response carries no payload (query {name}).')
                block.append(f'message {name}Response {{')
                block.append('  google.protobuf.Empty empty = 1;')
                block.append('}')
                block.append('')
            else:
                wm = re.search(r'WHERE\s+(?:\w+\.)?(\w+)\s*=\s*\$1', sql)
                fk_field = wm.group(1)
                pg_type = dict(cols)[fk_field]
                proto_type, _ = pg_type_to_proto(pg_type)
                if qtype == 'many':
                    block.append(f'// {name}Request asks for active dc.{table} rows filtered by {fk_field} (query {name}).')
                    block.append(f'message {name}Request {{')
                    block.append(f'  {proto_type} {fk_field} = 1;')
                    block.append('}')
                    block.append('')
                    block.append(f'// {name}Response returns active dc.{table} rows filtered by {fk_field} (query {name}).')
                    block.append(f'message {name}Response {{')
                    block.append(f'  repeated {entity} {e_plural} = 1;')
                    block.append('}')
                    block.append('')
                else:
                    block.append(f'// {name}Request asks for a single active dc.{table} row by {fk_field} (query {name}).')
                    block.append(f'message {name}Request {{')
                    block.append(f'  {proto_type} {fk_field} = 1;')
                    block.append('}')
                    block.append('')
                    block.append(f'// {name}Response returns a single active dc.{table} row (query {name}).')
                    block.append(f'message {name}Response {{')
                    block.append(f'  {entity} {e_snake} = 1;')
                    block.append('}')
                    block.append('')

        message_blocks.append('\n'.join(block).rstrip() + '\n')

    header = '''syntax = "proto3";

package user_domain_roles.v1;

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";

option go_package = "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1;userdomainrolesv1";

// UserDomainRolesService mirrors query.sql one-to-one:
// 8 standard CRUD calls per table.
service UserDomainRolesService {
'''
    service_body = '\n'.join(service_lines)
    out = header + service_body + '\n}\n\n' + '\n'.join(message_blocks)
    with open(OUT_PATH, 'w', encoding='utf-8', newline='\n') as f:
        f.write(out)
    print("written", OUT_PATH)


if __name__ == '__main__':
    build_proto()
