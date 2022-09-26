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

## Request & Response 

### Sign Up

- Descripción: registra el usuario.
- Path */signup*
- Method: POST

Request:
```bash 
curl --location --request POST 'http://localhost:5050/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "mayemail@myemail.com",
    "password": "123"
}'
```
Response:
```json 
{
    "id": "2FHVXHJlsEqgsnmpRYYTRJkISXU",
    "email": "mayemail@myemail.com"
}
```

### Login
- Descripción: loguea al usuario y retorna un token que se utiliza en el header para verificar el usuario. 
- Path */login*
- Method: POST

Request
```bash
curl --location --request POST 'http://localhost:5050/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "mayemail@myemail.com",
    "password": "123"
}'
```
Response 
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIyRkhWWEhKbHNFcWdzbm1wUllZVFJKa0lTWFUiLCJleHAiOjE2NjQzMjE5ODJ9.yijMTmkFBj55IaVhtZ8tFZkxKSy9KwecR2959ufrSro"
}
```

### Me
- Descripción: obtiene información sobre el usuario, por medio del token.
- Path */me*
- Method: GET

Request
```bash 
curl --location --request GET 'http://localhost:5050/me' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIyRkF5RWphVVpIaTUzRzMxeXZyQXdaaW5PR3AiLCJleHAiOjE2NjQxMjE4Nzd9.UtznSz25D4d7kSqA8G8LIO-TmamSpl5P1L_-dDEP51w'
```

Response
```json
{
    "id": "2FHVXHJlsEqgsnmpRYYTRJkISXU",
    "email": "mayemail@myemail.com",
    "password": ""
}
```

### Registrar un Post
- Descripión: registra un Post, valida el token.
- Path */post*
- Method: POST

Request
```bash 
curl --location --request POST 'http://localhost:5050/posts' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIyRkhWWEhKbHNFcWdzbm1wUllZVFJKa0lTWFUiLCJleHAiOjE2NjQzMjE5ODJ9.yijMTmkFBj55IaVhtZ8tFZkxKSy9KwecR2959ufrSro' \
--header 'Content-Type: application/json' \
--data-raw '{
    "post_content": "mi tercer post"
}'
```
Response
```json
{
    "id": "2FHWGQJYGz0v7QxZiMW0CAbVbIA",
    "post_content": "mi tercer post"
}
```

### Obtener un Post
- Descripción: obtiene una descripción por medio del id el post, valida e token.
- Path */post/:id*
- Method: GET

Request 
```bash
curl --location --request GET 'http://localhost:5050/posts/2FAyJ9G60hL0ZAMSZoDhD9dVFp5' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIyRkhWWEhKbHNFcWdzbm1wUllZVFJKa0lTWFUiLCJleHAiOjE2NjQzMjE5ODJ9.yijMTmkFBj55IaVhtZ8tFZkxKSy9KwecR2959ufrSro'
```

Response 
```json 
{
    "id": "",
    "post_content": "",
    "created_at": "0001-01-01T00:00:00Z",
    "user_id": ""
}
```

### Update Post
- Descripción: update del Post, valida el token y que el usuario que creo el Post genere esta acción.
- Path *posts/:id*
- Method: PUT

Request 
```bash 
curl --location --request PUT 'http://localhost:5050/posts/2FAyJ9G60hL0ZAMSZoDhD9dVFp5' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIyRkhWWEhKbHNFcWdzbm1wUllZVFJKa0lTWFUiLCJleHAiOjE2NjQzMjE5ODJ9.yijMTmkFBj55IaVhtZ8tFZkxKSy9KwecR2959ufrSro' \
--header 'Content-Type: application/json' \
--data-raw '{
    "post_content": "post modificado"
}'
```

Response 

```json 
{
    "message": "Post updated"
}
```

### Eliminar Post 
- Descripción: elimia un Post utilizando el id del Post y que el usuario quien realizo el Post pueda hacerlo, valida el token.
- Path */posts/:id*
- Method: DELETE

Request 
```bash 
curl --location --request DELETE 'http://localhost:5050/posts/2FAyJ9G60hL0ZAMSZoDhD9dVFp5' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIyRkhWWEhKbHNFcWdzbm1wUllZVFJKa0lTWFUiLCJleHAiOjE2NjQzMjE5ODJ9.yijMTmkFBj55IaVhtZ8tFZkxKSy9KwecR2959ufrSro' 
```

Response 
```json
{
    "message": "Post deleted"
}
```

### Paginar Posts
- Descripción: obtiene los Posts, valida el token. 
- Path */posts?page=:page*
- Method: GET

Request 
```bash
curl --location --request GET 'http://localhost:5050/posts?page=0' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIyRkgxbW9Zbm1mM01sQ2o5cXhhYUZiOTFWVmIiLCJleHAiOjE2NjQzMDcxNjB9.yEPTYBTxe5eIgOLtm48CL-Jl0OaAPDUQ5KROS1NAZTg'
```

Response
```json 
[
    {
        "id": "2FH1pKRTtASWk5WouiTZpTB8fkM",
        "post_content": "mi primer post",
        "created_at": "2022-09-25T19:32:56.01398Z",
        "user_id": "2FH1moYnmf3MlCj9qxaaFb91VVb"
    },
    {
        "id": "2FH1pv8MMj3gGMSvDV1HzAtNBTB",
        "post_content": "mi segundo post",
        "created_at": "2022-09-25T19:33:01.403763Z",
        "user_id": "2FH1moYnmf3MlCj9qxaaFb91VVb"
    },
    {
        "id": "2FH1qXpeFIUNowaxq2Y14PWuvA0",
        "post_content": "mi tercer post",
        "created_at": "2022-09-25T19:33:06.893616Z",
        "user_id": "2FH1moYnmf3MlCj9qxaaFb91VVb"
    }
]
```

### Websocket 
- Descripción: notifica por medio de websocket la creación de los diferentes Post
- Path */ws*
- Method: Websocket

Datos de conexión:
```bash 
Header -> Authorization: obtener el token del servicio login
host: localhost:5050/ws
```