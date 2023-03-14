# Merge Coding Assignment

## Requirements
- Ability to create account with two roles (admin, user) and log in.

- Admin should be able to:
	- Add items
	- Suspend user

- User should be able to:
	- List available items
	- Add items to a cart (if there are items in stock)
	- Remove items from their cart

- Restrict the access to APIs through RBAC mechanism.
  
## Schema

1. Accounts
ID, Name, Email, Password, Role (Enum), Active(Enum)

2. Items
ID, Name, SKU

3. Cart
ID, AccountID, ItemID, ItemName

## REST APIs

1. Admin
	1. __Add Item__
		Method: POST
		URI: /api/v1/admin/items
		Request Body: {'name': 'xyz', 'sku': 5}
        Headers:  {'Authorization': 'Bearer xxxxxxx', 'Content-Type': 'application/json'}
		Response Codes: 201, 400, 403, 500
		Content-Type: application/json
	2. __Suspend User__
		Method: PUT
		URI: /api/v1/accounts/suspend
        Headers:  {'Authorization': 'Bearer xxxxxxx', 'Content-Type': 'application/json'}
		Response Code: 200, 400, 403, 500
		Request Body: {'email': 'not@allowed.in'}
		
2. Users
	1. __List Available Items in store__
		Method: GET
		URI: /api/v1/user/items
        Headers:  {'Authorization': 'Bearer xxxxxxx', 'Content-Type': 'application/json'}
		Response Body:  [{'name': 'xyz', 'sku': 5}, ...]
		Response Codes: 200, 403, 500
	2. __Add Item to Cart__
		Method: POST
		URI: /api/v1/user/cart-items
        Headers:  {'Authorization': 'Bearer xxxxxxx', 'Content-Type': 'application/json'}
		Request Body: {'id': 1, 'name': 'xyz'}
		Response Codes: 200, 400, 403, 500
	3. __Remove Items from Cart__
		Method: DELETE
		Headers:  {'Authorization': 'Bearer xxxxxxx'}
		URI: /api/v1/user/cart-items?item_id={x}
		Response Codes: 204, 400, 403, 500
3. Accounts
	1. __Create account__
		Method: POST
		URI: /api/v1/accounts/register
		Request Body:  {'name': 'xss', 'email': 'dasd@sda.com', 'password': 'somethinglong', 'role': 'user'}
        > Note: allowed `role` values are `user` and `admin`.

		Response Codes: 201,  400, 500
	2. __Log In__
        Method: POST
		URI: /api/v1/accounts/login
		Request Body:  {'username': 'dasd@sda.com', 'password': 'rsewdsa'}
		Response Codes: 200,  400, 500
	3. __Log Out__
	    Method: POST
		URI: /api/v1/accounts/logout
		Headers:  {'Authorization': 'Bearer xxxxxxx'}
		Response Codes: 200,  400, 500

## Directory Structure and Code Organization

## System Requirements

1. Go version: 1.20
2. Postgres

## Setup Procecure

1. Create an `application.yaml` file in the root of the project directory.
2. Copy the following configurations in the `application.yaml` file:
    ```
    APP_PORT: '8888'
    JWT_REALM: 'test zone'
    JWT_SECRET: 'secret key'
    DATABASE_NAME: 'shop'
    DATABASE_HOST: 'localhost'
    DATABASE_PORT: '5432'
    DATABASE_USER: 'owner'
    DATABASE_PASSWORD: 'secret'
    DATABASE_DIALECT: 'postgres'
    DATABASE_SSL_MODE: 'disable'
    DATABASE_MIGRATIONS_DIR: 'internal/migrations'
    ```
3. Using any Postgres client: CLI app or a GUI, run the following commands:
    ```
    CREATE ROLE "owner" with login password 'secret';
    CREATE DATABASE shop;
    ```
4. From the root of the project directory, run the following command:
    ```
    go run cmd/cli/main.go db:migrate:up
    ```
    > Running the above command will create the necessary tables in the database.
6. To finally run the webserver, from the root of the project directory, run the following command:
    ```
    go run cmd/cli/main.go start:webserver
    ```
    
## Example Requests/Responses

> _Note:_ The examples are showcased using an HTTP client library called `httpie`
> To install `httpie`, run `brew install httpie` in the terminal.

1. __Register__
    Request:
    ```
    http -v --json POST "localhost:8888/api/v1/accounts/register" name=mayank email=mayank@admin.com password=secretpass role=admin  "Content-Type: application/json"
    ```
    __Example Response:__
    ```
    {"message": "account created successfully.",
    "status": "success"}
    ```
2. __Login__
    Request:
    ``` 
    http -v --json POST localhost:8888/api/v1/accounts/login email=mayank@admin.com password=secretpass
    ```
    __Example Response:__
    ```
    {"code": 200,
    "expire": "2023-03-13T11:21:51+05:30",
    "token": "xxxxxx"}
    ```
3. __Logout__
    __Request:__
    ```
    http -v --json POST localhost:8888/api/v1/accounts/logout "Authorization: Bearer xxxx"  "Content-Type: application/json"
    ```
    __Response:__
    ```
    {"code": 200}
    ```
4. __Suspend User__
    __Request:__
    ```
    http -v --json PUT "localhost:8888/api/v1/accounts/suspend" email=may@dju.com "Authorization: Bearer xxxx"
    ```
    __Response:__
    ```
    {"message": "account suspended successfully.",
    "status": "success"}
    ```
5. __List Items__
    __Request:__
    ``` 
    http -f GET "localhost:8888/api/v1/user/items" "Authorization: Bearer xxxxx"
    ```
    __Example Response:__
    ```
    [{
        "name": "chair",
        "sku": "5"
    },
    {
        "name": "table",
        "sku": "5"
    },
    {
        "name": "amp",
        "sku": "5"
    }]
    ```
6. __Add Item To Cart__
    __Request:__
    ```
    http -v --json POST "localhost:8888/api/v1/user/cart-items" item_id=1 "Authorization: Bearer xxxxxx"
    ```
    __Example Response:__
    ```
    {"message": "item added successfully",
    "status": "success"}
    ```
7. __Remove Item From Cart__
    __Request:__
    ```
    http DELETE "localhost:8888/api/v1/user/cart-items?item_id=1" "Authorization: Bearer xxxxx"
    ```
    > item_id query param refers to the ID of the item in the "items" table.
    
    __Example Response:__
    No response. Only 204 No Content response is returned
    
8. __Ping__
    __Request:__
    ```
    http GET "localhost:8888/api/ping"
    ```
    __Example Response:__
    ```
    {"ping": "pong"}
    ```
9. __Refresh Auth Token__
    __Request:__
    ```
    http GET "localhost:8888/api/v1/auth/refresh_token" "Authorization: Bearer xxxx"
    ```
    __Example Response:__
    
    ```
    {"code": 200,
    "expire": "2023-03-13T16:43:13+05:30",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiNyIsImV4cCI6MTY3ODcwNTk5MywiaWQiOiJtYXlAZGp1LmNvbSIsImlzX2FjdGl2ZSI6dHJ1ZSwib3JpZ19pYXQiOjE2Nzg3MDIzOTMsInJvbGUiOiJ1c2VyICJ9.h8x9IYXiw3ewCp-DQEnp4FU63o9Ceog6wJIA-XzDBgg"}

    ```
 
## TODOs (_or Things to Improve_):

1. Add better validation of request inputs.
2. Return "client-friendly" standard and error responses.
3. Add more unit tests.
4. Add more API tests.
5. Explore more edge cases.
6. I have implemented RBAC using JWT-based middleware which works fine when the business use-cases are limited, but would be a pain when the scope of the system increases. So, implementing RBAC using [Casbin](https://github.com/casbin/casbin) would be better for the long run.
7. Decrease the number of files in the `serializer` directory.
8. Improve route names.
9. Add contextual logging.
10. Add hashed password to DB.

## Notes

1. For the sake of focussing entirely on the problem statement, I have elided certain details. For example, there's no `price` attribute (among many other attributes) in the `items` schema.
2. I have used the DDD (Domain-Driven Design) approach to build the project, as this helps with testing each layer while only depending on the layer underneath it.
3. Directory Structure:
    1. `internal/api`: The HTTP API handlers exist here.
    2. `internal/domain`: The entities that reflect the nature of the business exist here.
    3. `internal/middleware`: All the HTTP middlewares exist here. (only one exists now)
    4. `internal/migrations`: Contains the SQL migration script and the code to run the migrations on Postgres.
    5. `internal/serializer`: All the Request and Response types exist here.
    6. `internal/repository`: Contains code that communicates with the database. It is used by the `service` layer.
    7. `internal/service`: Contains the code that communicates with the repository. 
    8. `pkg/database`: Contains the database-specific code.
    9. `cmd/cli`: Contains the code to execute the code.
    10. `config`: Contains code to load configurations from a YAML file.