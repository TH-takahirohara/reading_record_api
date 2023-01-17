#!/bin/sh

/wait
/migrate \
  -path $MIGRATIONS_DIR \
  -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" \
  $@
