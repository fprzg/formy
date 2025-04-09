#! /bin/bash

curl -X POST http://localhost:3000/api/submission/new/1 \
  -d "name=Penpals" \
  -d "email=pen@pals.com" \
  -d "subject=Consulta sobre servicios" \
  -d "message=Hola, me gustaría saber más sobre su oferta de diseño web."
