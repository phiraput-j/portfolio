# Mock API Project

This is a **mock API project** built in **Go** for practicing API testing. It supports **user and product management** with **token-based authentication** and **admin/user roles**. It also includes a **frontend HTML UI** for testing each API action.

---

## Features

- User registration and login.
- JWT-like tokens (base64-encoded JSON) with **1-minute expiry**.
- Role-based access: `admin` or `user`.
- Product CRUD (create, read, update, delete) with admin-only access for create/update/delete.
- In-memory storage (no database needed).
- Frontend UI (`index.html`) to test each API with:
  - Header inputs
  - JSON body input
  - Response output
- API documentation (`api-spec.html`) for QA reference.

---

## How to install go
1. Open PowerShell **as Administrator** and run:
2. Install Go via Chocolatey:
```powershell
choco install golang -y
```
3. Close and reopen PowerShell or Command Prompt, then verify installation:
```powershell
go version
```
You should see output like:
```bash
go version go1.21.2 windows/amd64
```
*Chocolatey automatically sets the PATH, so you can run go from any terminal.*


## Backend Setup (Go)

1. Clone the repository.
2. Run the server:

```bash
go run main.go
```

Server runs at: http://localhost:8080

* http://localhost:8080/index.html (Frontend UI)
* http://localhost:8080/api-spec.html (API documentation)

API Endpoints

### User APIs
| Method | Endpoint    | Description                                    | Role   |
| ------ | ----------- | ---------------------------------------------- | ------ |
| POST   | /register   | Register new user                              | Public |
| POST   | /login      | Login user and get token (expires in 1 minute) | Public |
| GET    | /users      | Get all users                                  | Admin  |
| GET    | /users/{id} | Get user by ID                                 | Admin  |

### Product APIs
| Method | Endpoint       | Description      | Role               |
| ------ | -------------- | ---------------- | ------------------ |
| GET    | /products      | Get all products | Any logged-in user |
| POST   | /products      | Create a product | Admin              |
| PUT    | /products/{id} | Update product   | Admin              |
| DELETE | /products/{id} | Delete product   | Admin              |

*Note: All protected endpoints require Authorization: Bearer <token> header. Tokens expire 1 minute after login.*

## How to Test

1. Register a new user:

```json
{
  "username": "alice",
  "password": "mypassword",
  "role": "admin"
}
```

2. Login with the user and get the token.

3. Use the token in Authorization header for protected endpoints.

4. Test user and product APIs in index.html UI.

## Token Expiry

* Login returns a token with expires_in: 60s.

* After 1 minute, all requests with that token will return 401 Token expired.

## Validation Rules

* User Registration: username must be unique; username/password required.

* Product: title, description cannot be empty; price must be positive.

* Admin-only endpoints: creating/updating/deleting products, getting all users or a user by ID.

## Notes

* This is a mock API; all data is in-memory. Restarting the server clears all users/products.

* Tokens are simple base64 JSON, not secure for production. Used for learning/testing.