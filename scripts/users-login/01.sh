#! /bin/bash

set -xe

curl -X POST http://localhost:3000/users/login -d "user_name=alice&password=securepass" -c cookies.txt &&
    cat cookies.txt &&
    rm ./cookies.txt