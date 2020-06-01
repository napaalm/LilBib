#!/bin/sh

docker exec database_db_1 sh -c 'exec mysqldump --all-databases -uroot -p"$MYSQL_ROOT_PASSWORD"' > test_database.sql
