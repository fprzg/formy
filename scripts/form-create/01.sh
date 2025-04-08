#! /bin/bash

curl -X POST http://localhost:3000/api/form/create \
  -d "user_id=1" \
  -d "name=Juan Pérez" \
  -d "description=Perfil de ejemplo" \
  -d "field_name=Email" \
  -d "field_type=string" \
  -d "field_constraints=[{\"constraint_name\": \"required\"},{\"constraint_name\": \"email\"}]" \
  -d "field_name=Age" \
  -d "field_type=int" \
  -d "field_constraints=[{\"constraint_name\": \"interval\", \"min\": 0, \"max\": 150}]"