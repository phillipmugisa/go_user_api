# go_user_api
Go implementation of an user account management api

- Install
  * Go
  * Docker
  * docker-compose

- Usage: 

```bash
    git clone github.com/phillipmugisa/go_user_api

    cd go_user_api

    docker-compose up
```

## Set up test database

```bash
    docker exec -it api_database bash -l

    mysql -u root -p
```

## Api endpoint inside the container:

### Creating User
- localhost:8000/user (POST)

#### Request Data
```json
{
  "username": "<string>",
  "email": "<string>",
  "password": "<string>",
  "region": "<string>",
  "userLanguage": "<string>",
  "userDateBirth": "<string>",
  "firstName": "<string>",
  "lastName": "<string>",
  "phone": "<string>",
  "userGender": "<string>"
}
```

#### Response Data
```json
{
  "username": "<string>",
  "code": "<string>"
}
```

### Verifying User using code sent to email
- localhost:8000/user/checkOtpcode (POST)

#### Request Data
```json
{
  "username": "<string>",
  "code": "<string>"
}
```

#### Response Data
```html
HTTP-code 200 OK
```

## Api endpoint out of container:
```bash
    go run .

    # using flag
    go run . -listenAddr "3000"
```


## Using make:
```bash
    make run
```
