# gicicm

## TO RUN 
```
To run the server: 
docker-compose up
```

## RUN TESTS 
```
for unit tests:
make unit 

for integration tests:
make integration

for all tests: 
make testall
```

## API's 
```
Create User: 
POST /gicicm/auth/signup HTTP/1.1
Host: localhost:8000
Content-Type: application/json
{
    "email":"test@gmail.com",
	"name":"clayton gonsalves",
	"password":"helld%Fo123"
}

Login: 
POST /gicicm/auth/login HTTP/1.1
Host: localhost:8000
Content-Type: application/json
{
    "email":"test@gmail.com",
	"password":"hello123"
}

Logout: 
POST /gicicm/auth/logout HTTP/1.1
Auth: Bearer type

List Users:
GET /gicicm/users HTTP/1.1
Host: localhost:8000
Auth: Bearer type

Delete User
DELETE /gicicm/users/{email} HTTP/1.1
Host: localhost:8000
Auth: Bearer type

```

## TODO's/ Improvements

-  separate error handler package for error handling 
- better security management
- use docker test for integration tests
- more ut/it coverage 





