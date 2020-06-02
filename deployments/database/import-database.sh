#!/bin/sh

docker exec -i database_db_1 sh -c 'exec mysql -uroot -p"secret"' < lilbib_example.sql
