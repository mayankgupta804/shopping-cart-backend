# Merge Coding Assignment

  

- Ability to create account with two roles (admin, user) and log in.

	- Use basic auth.

	- When creating a new account, if role is set to admin in the POST request body, then create an admin account.

- Admin should be able to

	- Add items

	- Suspend user

- User should be able to

	- List available items

	- Add items to a cart (if there are items in stock)

	- Remove items from their cart

- Restrict the access to APIs through RBAC mechanism

  

# Questions

  
1. Can I keep the *List Available Items* API simple - meaning no paginated items?

2. Also, is it fine if I go with a limited set of items, let's say three or four different items with 'x' SKU?

3. For the sake of simplicity, I am using a sqlite DB. Is the choice of DB something that I should consider for this assignment?

4. Can I make the login functionality simpler by just using Basic Auth for every request?

  
## Schema

1. Accounts
ID, Name, Email, Password, Role (Enum), Active(Enum)

2. Items
ID, Name, SKU

3. Cart
ID, AccountID, AccountName, AccountEmail, AccountActive(Enum)

5. CartItems
ID, AccountID, ItemID, ItemName


## REST APIs

1. Admin
	1. Add Items (Bulk addition)
		Method: POST
		URI: /items
		Request Body: [{'id': 1, 'name': 'xyz', 'sku': 5}, ...]
		Response Body: {'message': 'Items added to stock'}
		Response Codes: 201 Created, 400 Bad Request, 401, 403
		Content-Type: application/json
	2. Suspend User
		Method: PUT
		URI: /users/suspend
		Response Code: 200 OK, 400, 401, 403
		Request Body: {'id': 1}
		
2. Users
	1. List Available Items
		Method: GET
		URI: /items
		Response Body:  [{'id': 1, 'name': 'xyz', 'sku': 5}, ...]
		Response Codes: 200,  401
	2. Add Items to Cart
		Method: PATCH
		URI: /cart-items/add
		Request Body: [{'id': 1, 'name': 'xyz', 'count': 1}, ....]
		Response Codes: 200,  401, 400
	3. Remove Items from Cart
		Method: PATCH
		URI: /cart-items/remove
		Request Body: [{'id': 1, 'name': 'xyz'}, ....] OR {'id': 1, 'name': 'xyz'}
		Response Codes: 204 No Content, 400 Bad Request, 401
3. Account
	1. Create account
		Method: POST
		URI: /account
		Request Body:  {'name': 'xss', 'email': 'dasd@sda.com', 'password': 'dasddas'}
		Response Codes: 201,  400
	2. Log In