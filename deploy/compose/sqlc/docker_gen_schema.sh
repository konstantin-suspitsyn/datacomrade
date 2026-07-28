#!/bin/bash
export PGPASSWORD=$1

pg_dump -U $2 -s -d $3 > schema.sql