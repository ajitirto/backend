# API Testing Guide (HTTPie)

Base URL

```text
http://localhost:8080/health
```

---

# Authentication

## Sign Up

```bash
http POST :8080/api/authsignup \
    name="John Doe" \
    email="john@example.com" \
    password="password123"
```

---

## Login

```bash
http POST :8080/api/auth/login \
    email="john@example.com" \
    password="password123"
```

**Example Response**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

Simpan token:

```bash
export TOKEN="<JWT_TOKEN>"
```

---

# Profile

## Get Profile

```bash
http GET :8080/api/profile \
    Authorization:"Bearer $TOKEN"
```

---

# Product

## Create Product

```bash
http POST :8080/api/products \
    Authorization:"Bearer $TOKEN" \
    name="MacBook Pro" \
    price:=25000000 \
    stock:=10
```

---

## Get All Products

```bash
http GET :8080/api/products
```

---

## Get Product By ID

```bash
http GET :8080/api/products/1
```

---

## Update Product

```bash
http PUT :8080/api/products/1 \
    Authorization:"Bearer $TOKEN" \
    name="MacBook Pro M4" \
    price:=27000000 \
    stock:=15
```

---

## Delete Product

```bash
http DELETE :8080/api/products/1 \
    Authorization:"Bearer $TOKEN"
```

---

# Authentication Flow

```text
Sign Up
   │
   ▼
Login
   │
   ▼
Copy access_token
   │
   ▼
export TOKEN="<JWT_TOKEN>"
   │
   ▼
Profile
   │
   ├───────────────┐
   ▼               ▼
Create Product   Get Products
   │
   ▼
Update Product
   │
   ▼
Delete Product
```

---

# HTTP Status Codes

| Status | Description |
|--------|-------------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 500 | Internal Server Error |

---

# Tips

 Menyimpan token di variable bash

```bash
TOKEN=$(
  http --body POST :8080/api/auth/login \
    email=admin@example.com \
    password=password123 \
  | jq -r '.token'
)

echo "$TOKEN"
```
