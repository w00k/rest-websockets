# Dependences

* go get -u github.com/gorilla/mux
* go get -u github.com/gorilla/websocket
* go get -u github.com/golang-jwt/jwt
* go get -u github.com/joho/godotenv
* go get -u github.com/lib/pq
* go get -u github.com/segmentio/ksuid
* go get -u golang.org/x/crypto/bcrypt

## Patrón repositorio

Inversión de dependencias, los códigos deben estar basados en abstraciones y no en cosas comcretas.

// Abstracciones
Handler - GetUserByIdPostgres 
        - GetUserByIdMongoDB

// Concretas
Handler - GerUserById - return User
          Postgres
          MongoDB
          etc 

## Docker 

1.- Crear el contenedor 
```bash 
$ docker build . -t user-ws-rest-db 
```

2.- Levantar el contenedor
```bash
$ docker run -p 54321:5432 user-ws-rest-db
```

## Request 

Path */signup*
```json 
{
    "email": "mayemail@myemail.com",
    "password": "123"
}
```
