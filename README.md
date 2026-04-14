# CBT E-Commerce Platform

A modern, production-grade e-commerce application built with a **Go backend**, **React frontend**, **Supabase PostgreSQL**, and **Stripe payments**. This project demonstrates clean architecture, secure authentication, and real-world payment processing.

---

## 🚀 Project Overview

**CBT Store** is a full-stack e-commerce platform designed as a **technical interview assignment**. It showcases best practices in:

- **Clean Architecture:** Clear separation of concerns with Handler → Service → Repository → DB layers
- **Security:** JWT-based authentication, bcrypt password hashing, Stripe webhook signature verification
- **Scalability:** Transaction-based order creation, proper error handling, and type safety
- **User Experience:** Responsive Tailwind CSS design, smooth checkout flow, real-time cart state

### Key Capabilities

✅ Browse products from PostgreSQL database  
✅ Add items to cart with persistent localStorage state  
✅ Register/Login with JWT tokens  
✅ Checkout via Stripe Checkout (Test Mode)  
✅ Webhook-based order confirmation  
✅ Beautiful, minimalistic UI  

---

## 🛠 Tech Stack

### **Backend**
- **Language:** Go 1.21
- **Routing:** `go-chi/chi/v5` - Lightweight, composable router
- **Database:** `jackc/pgx/v5` - High-performance PostgreSQL driver
- **Authentication:** `golang-jwt/jwt/v5` - JWT token management
- **Security:** `golang.org/x/crypto/bcrypt` - Password hashing
- **Payments:** `stripe/stripe-go/v78` - Stripe API integration
- **Config:** `joho/godotenv` - Environment variable management

### **Frontend**
- **Framework:** React 18.2
- **Build Tool:** Vite 5.0
- **Routing:** `react-router-dom` 6.20
- **HTTP Client:** Axios 1.6
- **Styling:** Tailwind CSS 3.3
- **State Management:** React Hooks + Context API

### **Database**
- **Platform:** Supabase (PostgreSQL)
- **Tables:** Users, Products, Orders, OrderItems
- **Features:** UUID primary keys, timestamps, foreign key constraints, indexes

### **Payment Processing**
- **Provider:** Stripe (Test Mode)
- **Features:** Checkout Sessions, Webhook signature verification, Payment status tracking

### **Deployment Targets**
- **Backend:** Render.com or Fly.io (free tier Go hosting)
- **Frontend:** Vercel (optimized for Next.js/React)

---

## 🌟 Features Implemented

### Core Features
- ✅ **Product Catalog** - Fetched from Supabase with images, prices, stock
- ✅ **Shopping Cart** - Client-side state with localStorage persistence
- ✅ **User Authentication** - Secure register/login with JWT tokens (24h expiry)
- ✅ **Stripe Checkout** - Payment processing with client reference ID
- ✅ **Order Management** - Orders linked to users with line items
- ✅ **Success/Error Pages** - Beautiful post-payment user feedback

### Bonus Features (Advanced)
- ✅ **JWT Token Verification Middleware** - Protected routes require valid Bearer tokens
- ✅ **Stripe Webhook Handler** - Cryptographically verified `checkout.session.completed` events
- ✅ **Database Transactions** - Order + OrderItems created atomically
- ✅ **Password Hashing** - bcrypt with cost factor 10
- ✅ **CORS Middleware** - Frontend can make cross-origin requests
- ✅ **Product Seeding** - SQL script to populate database with mock data
- ✅ **Clean Error Handling** - Explicit error wrapping and meaningful HTTP responses

---

## 📦 Local Setup

### Prerequisites
- **Go 1.21+** ([Download](https://golang.org/dl/))
- **Node.js 18+** ([Download](https://nodejs.org/))
- **Supabase Account** ([Free tier](https://supabase.com/))
- **Stripe Account** ([Free tier](https://stripe.com/))
- **Stripe CLI** (optional, for webhook testing)

---

### Backend Setup

#### 1. Navigate to Backend Directory
```bash
cd backend
```

#### 2. Create `.env` File
```bash
# Copy the template below and fill in your credentials
cat > .env << 'EOF'
DATABASE_URL=postgresql://[user]:[password]@[host]:5432/[database]
PORT=8080
JWT_SECRET=your_super_secret_jwt_key_change_in_production
STRIPE_SECRET_KEY=sk_test_your_actual_key_here
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret_here
EOF
```

**Get your DATABASE_URL from Supabase:**
1. Go to Supabase Dashboard → Project Settings → Database
2. Copy the connection string (URI)
3. Replace placeholders with your credentials

**Get Stripe keys from [Stripe Dashboard](https://dashboard.stripe.com):**
1. Go to Developers → API Keys
2. Copy Secret Key (starts with `sk_test_`)
3. For webhooks, go to Webhooks → Add endpoint

#### 3. Download Dependencies
```bash
go mod download
go mod tidy
```

#### 4. Run the Server
```bash
go run cmd/server/main.go
```

You should see:
```
✓ Database connection established
🚀 Server starting on port 8080
```

**Test the API:**
```bash
curl http://localhost:8080/health
# Expected: {"status": "ok"}
```

---

### Frontend Setup

#### 1. Navigate to Frontend Directory
```bash
cd frontend
```

#### 2. Create `.env.local` File
```bash
cat > .env.local << 'EOF'
VITE_API_URL=http://localhost:8080/api
VITE_STRIPE_PUBLIC_KEY=pk_test_your_actual_key_here
EOF
```

**Get your STRIPE_PUBLIC_KEY from [Stripe Dashboard](https://dashboard.stripe.com):**
1. Go to Developers → API Keys
2. Copy Publishable Key (starts with `pk_test_`)

#### 3. Install Dependencies
```bash
npm install
```

#### 4. Start Development Server
```bash
npm run dev
```

The app will open at `http://localhost:5173`.

---

### Stripe CLI Setup (Optional, for Webhook Testing)

The Stripe CLI allows you to test webhooks locally without deployment.

#### 1. Install Stripe CLI
**macOS:**
```bash
brew install stripe/stripe-cli/stripe
```

**Windows (PowerShell):**
```powershell
choco install stripe
```

**Linux:**
```bash
curl https://files.stripe.com/stripe-cli/releases/latest/linux/x86_64/stripe_linux_windows_amd64.tar.gz -o stripe_linux_windows_amd64.tar.gz
tar -xf stripe_linux_windows_amd64.tar.gz
sudo mv stripe /usr/local/bin
```

#### 2. Login and Forward Webhooks
```bash
stripe login
```

#### 3. Forward Events to Your Local Backend
```bash
stripe listen --forward-to localhost:8080/api/webhook/stripe
```

You'll get a webhook signing secret. Add it to `backend/.env`:
```env
STRIPE_WEBHOOK_SECRET=whsec_xxxxx
```

#### 4. Trigger Test Events
```bash
stripe trigger payment_intent.succeeded
stripe trigger checkout.session.completed
```

---

## 🔑 Environment Variables Template

### Backend `.env`
```env
# Database
DATABASE_URL=postgresql://user:password@host:5432/database

# Server
PORT=8080

# JWT Authentication
JWT_SECRET=your_super_secret_jwt_key_change_in_production

# Stripe API
STRIPE_SECRET_KEY=sk_test_xxxxx
STRIPE_WEBHOOK_SECRET=whsec_xxxxx
```

### Frontend `.env.local`
```env
# Backend API
VITE_API_URL=http://localhost:8080/api

# Stripe Public Key
VITE_STRIPE_PUBLIC_KEY=pk_test_xxxxx
```

---

## 📡 API Endpoints

### Public Routes
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/products` | Fetch all products |
| `POST` | `/api/auth/register` | Register new user |
| `POST` | `/api/auth/login` | Login user |
| `POST` | `/api/webhook/stripe` | Stripe webhook listener |
| `GET` | `/health` | Health check |

### Protected Routes (Require JWT)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/checkout` | Init Stripe Checkout Session |

---

## 🗂 Project Structure

```
CBT/
├── backend/
│   ├── cmd/server/
│   │   └── main.go              # Entry point
│   ├── internal/
│   │   ├── database/
│   │   │   └── db.go            # DB connection
│   │   ├── handlers/
│   │   │   ├── auth.go          # Register/Login
│   │   │   ├── products.go      # Product endpoints
│   │   │   ├── orders.go        # Checkout endpoint
│   │   │   └── webhook.go       # Stripe webhooks
│   │   ├── middleware/
│   │   │   ├── cors.go          # CORS headers
│   │   │   └── auth.go          # JWT verification
│   │   ├── models/
│   │   │   ├── user.go
│   │   │   ├── product.go
│   │   │   └── order.go
│   │   ├── repository/
│   │   │   ├── user.go
│   │   │   ├── product.go       # DB queries
│   │   │   └── order.go         # Order DB ops
│   │   └── services/
│   │       ├── auth.go          # Auth logic
│   │       ├── product.go       # Product logic
│   │       └── stripe.go        # Stripe integration
│   ├── scripts/
│   │   └── seed.sql             # Product seeding
│   ├── .env                     # (Add to .gitignore)
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Header.jsx       # Navigation
│   │   │   ├── ProductList.jsx  # Product grid
│   │   │   └── Cart.jsx         # Shopping cart
│   │   ├── pages/
│   │   │   ├── Home.jsx         # Main page
│   │   │   ├── Login.jsx        # Login form
│   │   │   ├── Register.jsx     # Registration form
│   │   │   ├── Success.jsx      # Post-payment
│   │   │   └── Cancelled.jsx    # Payment cancelled
│   │   ├── hooks/
│   │   │   ├── useCart.js       # Cart state
│   │   │   └── useAuth.jsx      # Auth context
│   │   ├── services/
│   │   │   └── api.js           # Axios instance
│   │   ├── App.jsx              # Router setup
│   │   ├── main.jsx             # React entry
│   │   └── index.css            # Tailwind
│   ├── .env.local               # (Add to .gitignore)
│   ├── vite.config.js
│   ├── tailwind.config.js
│   ├── package.json
│   └── index.html
│
├── .gitignore
└── README.md
```

---

## 🧪 Testing the Full Flow

### 1. Register a New Account
```bash
POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}
```

### 2. Login
```bash
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}
```
Response includes `token` — save this for subsequent requests.

### 3. Get Products
```bash
GET http://localhost:8080/api/products
```

### 4. Initiate Checkout
```bash
POST http://localhost:8080/api/checkout
Authorization: Bearer <your_token>
Content-Type: application/json

{
  "items": [
    {
      "product_id": "...",
      "quantity": 2,
      "price": 149.99,
      "name": "Wireless Headphones"
    }
  ],
  "success_url": "http://localhost:5173/success",
  "cancel_url": "http://localhost:5173/cancelled"
}
```

### 5. Use Stripe Test Card
In Stripe Checkout, use:
- **Card:** `4242 4242 4242 4242`
- **Expiry:** Any future date (e.g., `12/25`)
- **CVC:** Any 3 digits (e.g., `123`)
- **Zip:** Any value

---

## 🚀 Deployment

### Deploy Backend (Render.com)
1. Push code to GitHub
2. Sign up at [Render.com](https://render.com)
3. Create new Web Service → Connect GitHub repo
4. Set environment variables in dashboard
5. Deploy!

### Deploy Frontend (Vercel)
1. Push code to GitHub
2. Sign up at [Vercel.com](https://vercel.com)
3. Import project → Connect GitHub
4. Add environment variables
5. Deploy!

---

## 📝 What I Learned

Building this project taught me invaluable lessons about **full-stack development** and **production-grade engineering**:

### Stripe Webhook Integration
Integrating Stripe webhooks was a humbling experience. The critical lesson: **always verify the cryptographic signature before processing webhook events**. Initially, I tried decoding the JSON directly, but Stripe requires the raw request body for signature verification. This taught me the importance of reading library documentation carefully and understanding security practices like HMAC verification. It's a great example of why simple "decode the JSON and process it" approaches fail in real-world scenarios.

### Go's Strict Typing and Imports
Go's strict type system and mandatory import usage (unused imports cause compilation failure) felt restrictive at first, but I now appreciate it as a feature. The language forces you to be intentional about your code. No silent type coercions, no dead imports cluttering your file. This discipline prevents entire classes of bugs that plague dynamically-typed languages. Managing pgx transactions, JSON unmarshaling, and interface{} types taught me to think about data flow and type safety from the ground up.

### Clean Architecture in Practice
Separating concerns into Handler → Service → Repository layers wasn't just theory. It made testing easier (mock repositories), maintenance simpler (logic is centralized), and scaling straightforward (each layer has a single responsibility). When I needed to change how orders are created, I only had to modify the repository and service—handlers stayed the same.

### Database Transactions Matter
Using SQL transactions to ensure Order + OrderItems are created atomically was eye-opening. If the order creates but items fail, the entire transaction rolls back. This prevents data corruption and ensures consistency—a critical requirement in payment systems.

### React Context for State Management
Building the `useAuth` hook with React Context and localStorage taught me that you don't always need Redux. For small-to-medium apps, Context + Hooks is elegant and sufficient. The key is thinking about state updates carefully (using callbacks to avoid stale closures).

### Tailwind CSS Productivity
Tailwind's utility-first approach initially seemed verbose, but the speed of building responsive UIs is incredible. No more context-switching between HTML and CSS files. The design system is consistent (spacing, colors, shadows) across the app without manually writing classes.

### Security Isn't Optional
From bcrypt password hashing to JWT token verification to webhook signature validation—security has to be baked into every layer. There's no "we'll add it later." Every endpoint, every user interaction, every external webhook must be treated as potentially hostile.

---

## 📚 Resources

- [Go Best Practices](https://golang.org/doc/effective_go)
- [Chi Router Documentation](https://github.com/go-chi/chi)
- [pgx/v5 Documentation](https://github.com/jackc/pgx)
- [Stripe API Reference](https://stripe.com/docs/api)
- [React Hooks Documentation](https://react.dev/reference/react)
- [Tailwind CSS Docs](https://tailwindcss.com/docs)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8949)

---

## 📄 License

This project is provided as-is for learning and interview purposes. Feel free to use, modify, and learn from it.

---

## 🙏 Acknowledgments

This project demonstrates modern web development practices, from secure authentication to real-world payment processing. Built with clean architecture, type safety, and a focus on user experience.

**Built with ❤️ for a technical interview challenge.**
