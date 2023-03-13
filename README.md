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
5. Running the above command will create the necessary tables in the database.
6. To finally run the webserver, from the root of the project directory, run the following command:
    ```
    go run cmd/cli/main.go start:webserver
    ```
    
## Example Requests/Responses

> _Note:_ The examples are showcased using an HTTP client library called `httpie`
> To install `httpie`, run `brew install httpie` in the terminal.

1. __Register__
2. __Login__
    Request:
    ``` 
    http -v --json POST localhost:8888/api/v1/accounts/login email=user@test.com password=something
    ```
    Response:
    ```
    {"code": 200,
    "expire": "2023-03-13T11:21:51+05:30",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiNyIsImV4cCI6MTY3ODY4NjcxMSwiaWQiOiJtYXlAZGp1LmNvbSIsImlzX2FjdGl2ZSI6dHJ1ZSwib3JpZ19pYXQiOjE2Nzg2ODMxMTEsInJvbGUiOiJ1c2VyICJ9.vUpbwXkpc_KfRAgHI1OCyqsYicHd78R_er_-dMpVWzQ"}
    ```
3. __Logout__
4. __Suspend User__
5. __List Items__
6. __Add Item To Cart__
7. __Remove Item From Cart__
8. __Ping__
9. __Refresh Auth Token__
 
## TODOs (_or Things to Improve_):
