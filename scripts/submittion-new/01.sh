#! /bin/bash

curl -X POST http://localhost:3000/api/submission/new/1 \
  -d "name=Lucía Fernanda" \
  -d "description=Perfil de ejemplo" \
  -d "email=lucifer@gmail.com"
