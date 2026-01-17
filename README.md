# E-Commerce Fullstack Monorepo

A premium, high-performance e-commerce platform built with a modern tech stack. This repository uses a **Monorepo** architecture to manage both the backend services and the web frontend in a single, unified codebase.

## 🚀 Tech Stack

### Backend (`/ecommerce-backend`)
- **Language:** Go (Golang)
- **Framework:** Gin Web Framework
- **Database:** PostgreSQL
- **ORM/Query:** Standard library / sqlx
- **Features:** JWT Authentication, RESTful API, Database Migrations.

### Frontend (`/ecommerce-web`)
- **Framework:** Next.js 15+ (App Router)
- **Styling:** Tailwind CSS
- **Language:** TypeScript
- **UI Components:** Shadcn UI (Radix UI)
- **State Management:** React Hooks & Context API
- **Features:** Responsive Design, Server-Side Rendering (SSR).

### Infrastructure
- **Orchestration:** Docker Compose
- **Database:** PostgreSQL 16 (Alpine)

---

## 📂 Project Structure

```text
.
├── ecommerce-backend/    # Go Backend API
├── ecommerce-web/        # Next.js Frontend
├── docker-compose.yml    # Docker orchestration for all services
└── .gitignore            # Root gitignore for the entire project
```

---

## 🛠️ Getting Started

### Prerequisites
- Docker & Docker Compose
- Go (for local development)
- Node.js & npm/yarn (for local development)

### Running with Docker (Recommended)
You can start the entire stack (Database, Backend, and Frontend) with a single command:

```bash
docker-compose up --build
```

- **Frontend:** [http://localhost:3000](http://localhost:3000)
- **Backend API:** [http://localhost:8080](http://localhost:8080)
- **Postgres:** `localhost:5432`

---

## 🔧 Development

### Backend Setup
1. Navigate to directory: `cd ecommerce-backend`
2. Configuration: Copy `.env.example` to `.env`
3. Run locally: `go run cmd/api/main.go`

### Frontend Setup
1. Navigate to directory: `cd ecommerce-web`
2. Install dependencies: `npm install`
3. Run development server: `npm run dev`

---

## 📝 Features Roadmap
- [x] Monorepo Architecture Setup
- [x] User Authentication (JWT)
- [x] Product Catalog & Search
- [ ] Shopping Cart Implementation
- [ ] Checkout Process & Payment Integration
- [ ] Mobile App Development (Coming Soon)

## 👤 Author
- GitHub: [@synchhans](https://github.com/synchhans)
