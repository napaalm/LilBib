#!/bin/sh

docker exec -i database_db_1 sh -c 'exec mysql -uroot -psecret' < lilbib_example.sql
