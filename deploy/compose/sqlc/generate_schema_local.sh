#!/bin/bash

# При запуске скрипта нужно передать путь к файлу .env

export $(grep -v '^#' $1 | xargs)

docker cp ./docker_gen_schema.sh $DATACATALOGUE_CONTAINER_NAME:docker_gen_schema.sh

docker exec $DATACATALOGUE_CONTAINER_NAME /bin/sh ./docker_gen_schema.sh $DATACATALOGUE_POSTGRES_PASSWORD $DATACATALOGUE_POSTGRES_USER $DATACATALOGUE_POSTGRES_DB

docker cp $DATACATALOGUE_CONTAINER_NAME:./schema.sql $2
