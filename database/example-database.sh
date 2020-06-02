#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# restart the database
! [ -z `docker-compose -f $DIR/docker-compose.yml ps -q db` ] || [ -z `docker ps -q --no-trunc | grep $(docker-compose -f $DIR/docker-compose.yml ps -q db)` ] && docker-compose -f $DIR/docker-compose.yml down
docker-compose -f $DIR/docker-compose.yml up -d

# import data (hacky hack)
printf "Waiting until MariaDB is fully started to create the database structure..."
until docker exec -i database_db_1 sh -c 'exec mysql -uroot -p"secret"' < $DIR/lilbib_example.sql 2> /dev/null
do
  sleep 1
done
echo " done!"
