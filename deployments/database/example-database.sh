#!/bin/sh

# restart the database
docker-compose down
docker-compose up -d

# import data (hacky hack)
until docker exec -i database_db_1 sh -c 'exec mysql -uroot -p"secret"' < lilbib_example.sql
do
  sleep 1
done
