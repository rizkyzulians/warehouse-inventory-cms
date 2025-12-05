# Warehouse Inventory Management System

A complete warehouse inventory management system built with Golang backend (REST API) and Next.js frontend with Tailwind CSS.

## 🚀 Features

### Backend (Golang)
- ✅ Clean Architecture with Repository Pattern
- ✅ JWT Authentication with Role-Based Access Control (Admin & Staff)
- ✅ RESTful API with standardized JSON responses
- ✅ Database Transactions for Purchase & Sales
- ✅ Automatic stock management
- ✅ Stock history tracking
- ✅ PostgreSQL database

### Frontend (Next.js + Tailwind CSS)
- ✅ Responsive dashboard
- ✅ Master Barang (Product Management)
- ✅ Stock Management
- ✅ Purchase Transactions
- ✅ Sales Transactions
- ✅ Stock History
- ✅ Protected routes based on user role
- ✅ Beautiful UI with Tailwind CSS

## 📋 Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- PostgreSQL 15 or higher
- Docker & Docker Compose (optional)

## 🛠️ Installation

### Option 1: Using Docker (Recommended)

1. **Clone the repository**
```bash
cd warehouse-inventory-cms
```

2. **Start all services**
```bash
docker-compose up -d
```

3. **Access the application**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Database: localhost:5432

4. **Login credentials**
- Admin: `admin` / `password123`
- Staff: `staff1` / `password123`

### Option 2: Manual Installation

#### Backend Setup

1. **Navigate to backend directory**
```bash
cd warehouse-api
```

2. **Install Go dependencies**
```bash
go mod download
```

3. **Setup environment variables**
```bash
copy .env.example .env
```

Edit `.env` with your database credentials:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=warehouse_db
JWT_SECRET=your-secret-key-here
PORT=8080
```

4. **Create database and run migrations**
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE warehouse_db;

# Exit psql
\q

# Run migrations
psql -U postgres -d warehouse_db -f migrations/001_create_tables.sql
psql -U postgres -d warehouse_db -f migrations/002_seed_data.sql
```

5. **Generate password hashes for seed users**

Create a small Go program to generate bcrypt hashes:
```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "password123"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    fmt.Println(string(hash))
}
```

Run it and update the hashed passwords in `migrations/002_seed_data.sql`

6. **Run the backend**
```bash
go run main.go
```

Backend will start on http://localhost:8080

#### Frontend Setup

1. **Navigate to frontend directory**
```bash
cd warehouse-frontend
```

2. **Install dependencies**
```bash
npm install
```

3. **Setup environment variables**
```bash
copy .env.local.example .env.local
```

Edit `.env.local`:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

4. **Run the frontend**
```bash
npm run dev
```

Frontend will start on http://localhost:3000

## 📚 API Documentation

### Authentication

#### Login
```http
POST /api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "jwt-token-here",
    "user": {
      "id": 1,
      "username": "admin",
      "nama": "Administrator",
      "role": "admin"
    }
  }
}
```

### Master Barang

All endpoints require `Authorization: Bearer <token>` header.

#### Get All Barang
```http
GET /api/barang?page=1&limit=10&search=laptop
```

#### Get Barang with Stock
```http
GET /api/barang/stok?page=1&limit=10&search=
```

#### Get Barang by ID
```http
GET /api/barang/{id}
```

#### Create Barang (Admin Only)
```http
POST /api/barang
Content-Type: application/json

{
  "kode_barang": "BRG011",
  "nama_barang": "Product Name",
  "kategori": "Electronics",
  "satuan": "Unit",
  "harga_beli": 100000,
  "harga_jual": 150000
}
```

#### Update Barang (Admin Only)
```http
PUT /api/barang/{id}
Content-Type: application/json

{
  "nama_barang": "Updated Product Name",
  "kategori": "Electronics",
  "satuan": "Unit",
  "harga_beli": 110000,
  "harga_jual": 160000
}
```

### Stock Management

#### Get All Stock
```http
GET /api/stok
```

#### Get Stock by Barang ID
```http
GET /api/stok/{barang_id}
```

#### Get Stock History
```http
GET /api/history-stok?page=1&limit=10
```

#### Get Stock History by Barang ID
```http
GET /api/history-stok/{barang_id}?page=1&limit=10
```

### Purchase (Pembelian)

#### Create Purchase
```http
POST /api/pembelian
Content-Type: application/json

{
  "no_faktur": "PO-2025-003",
  "tanggal": "2025-12-05",
  "supplier": "PT Supplier Example",
  "keterangan": "Purchase note",
  "details": [
    {
      "barang_id": 1,
      "qty": 5,
      "harga": 100000
    },
    {
      "barang_id": 2,
      "qty": 10,
      "harga": 50000
    }
  ]
}
```

**Business Logic:**
- Validates all barang exist
- Calculates subtotal and total automatically
- Updates stock (stok_akhir + qty)
- Inserts history_stok with jenis_transaksi = "masuk"
- Uses database transaction (rollback on error)

#### Get All Purchases
```http
GET /api/pembelian?page=1&limit=10
```

#### Get Purchase by ID
```http
GET /api/pembelian/{id}
```

### Sales (Penjualan)

#### Create Sale
```http
POST /api/penjualan
Content-Type: application/json

{
  "no_faktur": "SO-2025-003",
  "tanggal": "2025-12-05",
  "customer": "PT Customer Example",
  "keterangan": "Sale note",
  "details": [
    {
      "barang_id": 1,
      "qty": 2,
      "harga": 150000
    }
  ]
}
```

**Business Logic:**
- Validates all barang exist
- **Checks if stock is sufficient** (returns 400 with code "INSUFFICIENT_STOCK" if not)
- Calculates subtotal and total automatically
- Updates stock (stok_akhir - qty)
- Inserts history_stok with jenis_transaksi = "keluar"
- Uses database transaction (rollback on error)

#### Get All Sales
```http
GET /api/penjualan?page=1&limit=10
```

#### Get Sale by ID
```http
GET /api/penjualan/{id}
```

## 📊 Database Schema

### Tables

1. **users** - User authentication and roles
2. **master_barang** - Product master data
3. **mstok** - Stock information
4. **history_stok** - Stock movement history
5. **beli_header** - Purchase header
6. **beli_detail** - Purchase details
7. **jual_header** - Sales header
8. **jual_detail** - Sales details

See `warehouse-api/migrations/001_create_tables.sql` for complete schema.

## 🔐 Role-Based Access Control

### Admin Role
- Full access to all endpoints
- Can create, update, and delete master barang
- Can create purchases and sales
- Can view all data

### Staff Role
- Read access to master barang
- Can create purchases and sales
- Can view all data
- Cannot modify master barang

## 🧪 Testing

### Positive Test Cases

1. **Create Purchase** → Stock increases
2. **Create Sale** → Stock decreases
3. **Pagination** works correctly
4. **Stock history** appears after transactions

### Negative Test Cases

1. **Insufficient Stock** → Returns 400 with code "INSUFFICIENT_STOCK"
2. **Invalid Input** → Returns 422
3. **Item Not Found** → Returns 404
4. **Unauthorized** → Returns 401
5. **Forbidden (insufficient permissions)** → Returns 403

## 📁 Project Structure

```
warehouse-inventory-cms/
├── warehouse-api/                  # Backend (Golang)
│   ├── config/                     # Database configuration
│   ├── models/                     # Data models
│   ├── repositories/               # Repository layer
│   ├── services/                   # Business logic layer
│   ├── handlers/                   # HTTP handlers
│   ├── middleware/                 # JWT auth middleware
│   ├── migrations/                 # SQL migrations
│   ├── main.go                     # Entry point
│   ├── go.mod                      # Go dependencies
│   ├── Dockerfile                  # Docker configuration
│   └── .env.example                # Environment template
│
├── warehouse-frontend/             # Frontend (Next.js)
│   ├── src/
│   │   ├── components/             # Reusable components
│   │   ├── lib/                    # API client & auth utils
│   │   ├── pages/                  # Next.js pages
│   │   ├── styles/                 # Global styles
│   │   └── types/                  # TypeScript types
│   ├── public/                     # Static assets
│   ├── package.json                # Node dependencies
│   ├── Dockerfile                  # Docker configuration
│   ├── next.config.js              # Next.js configuration
│   ├── tailwind.config.js          # Tailwind configuration
│   └── tsconfig.json               # TypeScript configuration
│
├── docker-compose.yml              # Docker Compose configuration
└── README.md                       # This file
```

## 🚨 Error Response Format

All errors follow this format:

```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error description",
  "code": "ERROR_CODE"
}
```

### Error Codes

- `INSUFFICIENT_STOCK` - Not enough stock for sale (400)
- `VALIDATION_ERROR` - Invalid input (422)
- `NOT_FOUND` - Resource not found (404)
- `UNAUTHORIZED` - Authentication required (401)
- `FORBIDDEN` - Insufficient permissions (403)
- `INTERNAL_ERROR` - Server error (500)

## 🔧 Environment Variables

### Backend (.env)
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=warehouse_db
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 📦 Deployment

### Using Docker Compose

```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Manual Deployment

1. Build backend:
```bash
cd warehouse-api
go build -o warehouse-api
```

2. Build frontend:
```bash
cd warehouse-frontend
npm run build
npm start
```

## 🐛 Troubleshooting

### Backend won't start
- Check if PostgreSQL is running
- Verify database credentials in `.env`
- Ensure migrations have been run

### Frontend won't connect to backend
- Check `NEXT_PUBLIC_API_URL` in `.env.local`
- Verify backend is running on port 8080
- Check CORS settings in backend

### Docker issues
- Run `docker-compose down -v` to clean volumes
- Check Docker logs: `docker-compose logs`
- Ensure ports 3000, 8080, 5432 are not in use

## 📝 License

This project is created for educational purposes.

## 👥 Default Users

| Username | Password | Role |
|----------|----------|------|
| admin | password123 | admin |
| staff1 | password123 | staff |
| staff2 | password123 | staff |

## 🎯 Features Checklist

- [x] Clean Architecture
- [x] Repository Pattern
- [x] JWT Authentication
- [x] Role-Based Access Control
- [x] Database Transactions
- [x] Automatic Stock Management
- [x] Stock History
- [x] Standardized API Responses
- [x] Error Handling
- [x] Pagination
- [x] Search Functionality
- [x] Docker Support
- [x] Responsive UI
- [x] Protected Routes

## 🔮 Future Enhancements

- [ ] Unit Tests
- [ ] Integration Tests
- [ ] API Documentation (Swagger)
- [ ] Logging System
- [ ] Request Validation
- [ ] Rate Limiting
- [ ] Export to Excel/PDF
- [ ] Advanced Reporting
- [ ] Email Notifications
- [ ] Audit Trail

---

Built using Golang, Next.js, and Tailwind CSS