#! /bin/bash

curl -X POST http://localhost:3000/api/form/create \
  -d "user_id=123" \
  -d "name=Juan Pérez" \
  -d "description=Perfil de ejemplo" \
  -d "field_name=Age" \
  -d "field_type=int" \
  -d "field_constraints=min=0,max=150" \
  -d "field_name=Email" \
  -d "field_type=string" \
  -d "field_constraints=required,email"
