#! /usr/bin/bash

GREEN='\033[1;32m'
RESET='\033[0m'

#sudo apt -qq update > /dev/null
#sudo apt install jq curl python3 -y > /dev/null

echo -e "${GREEN}Phase 1 - Building images${RESET}"

docker stop grest mysql > /dev/null
docker rmi -f grest-testing mysql-testing

docker build -f src/test/dockerfile-grest -t grest-testing .
docker build -f src/test/dockerfile-mysql -t mysql-testing .


# Getting mysql container and database done
echo -e "${GREEN}Phase 2 - Starting MySQL${RESET}"

docker run --name mysql --rm --health-cmd='mysqladmin ping --silent' -d mysql-testing

while STATUS=$(docker inspect mysql --format "{{.State.Health.Status}}" $1); [ $STATUS != "healthy" ]; do

    sleep 1
  done


docker exec -i mysql sh -c 'exec mysql -h 127.0.0.1 -P 3306 -uroot -proot grest' < src/test/mysql.sql

# Getting GREST container done
echo -e "${GREEN}Phase 3 - Starting GREST${RESET}"

docker run --name grest --rm -d grest-testing ./grest


MYSQL_IP=$(docker inspect mysql | grep -Po -m 1 '(?<="IPAddress": )"\d+\.\d+.\d+\.\d+"' | sed 's/\"//g')
GREST_IP=$(docker inspect grest | grep -Po -m 1 '(?<="IPAddress": )"\d+\.\d+.\d+\.\d+"' | sed 's/\"//g')

sleep 5

echo -e "${GREEN}Phase 4 - Starting tests${RESET}"
python3 src/test/tester.py --grest $GREST_IP --mysql $MYSQL_IP

echo -e "${GREEN}Phase 5 - Cleaning things${RESET}"

docker stop mysql grest > /dev/null

echo "Cleaned"