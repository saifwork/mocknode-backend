# 🧩 Mock Service (MockNode Backend)

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![MongoDB](https://img.shields.io/badge/Database-MongoDB-47A248?logo=mongodb&logoColor=white)](https://www.mongodb.com)
[![Gin](https://img.shields.io/badge/Framework-Gin-blue?logo=go&logoColor=white)](https://gin-gonic.com)
[![Redis](https://img.shields.io/badge/Cache-Redis-DC382D?logo=redis&logoColor=white)](https://redis.io)

A modern backend built in **Golang** that allows developers to easily **create, manage, and test APIs** through projects and collections — the engine behind **MockNode** 🧠.  
It’s perfect for developers, learners, and testers who want to simulate API endpoints without building a custom backend.

---

## 🧭 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Tech Stack](#-tech-stack)
- [Setup Instructions](#-setup-instructions)
- [Environment Variables](#-environment-variables)
- [API Endpoints](#-api-endpoints)
- [Developer Notes](#-developer-notes)
- [Author](#-author)
- [License](#-license)

---

## 🚀 Features

### 🔐 Authentication
- Signup with email verification
- JWT-based login
- Forgot/Reset password (email link)
- Change password (authenticated users)
- Auto email sending via Gmail SMTP

### 📦 Projects
- Create, update, delete projects
- View all projects for a user
- Free users can create up to **2 projects**

### 📚 Collections
- Add dynamic fields to collections
- Field definitions include:
  - `string`, `number`, `boolean`, `email`, `date`, `array`, `object`, etc.
  - Constraints: `min`, `max`, `required`, `pattern`
- Free users can create up to **3 collections per project**

### 🗂️ Records
- CRUD operations for dynamic records
- Data validation based on collection schema
- Auto-updated timestamps

### ⚙️ Config
- `/config/types` — List supported field types
- `/config/presets` — Prebuilt mock collections (`users`, `products`, `posts`, `todos`, `comments`, `carts`)
- No authentication required for preset data

### 🧰 Health Check
Quick endpoint to check service status:


### Architecture ---

+-----------------------------------------------------+
|                   MockNode Backend                  |
|-----------------------------------------------------|
|  AuthService      -> Manages users, login, JWTs     |
|  ProjectService   -> User projects CRUD              |
|  CollectionService-> Schema & validation management  |
|  RecordService    -> Dynamic data CRUD               |
|  PresetService    -> Prebuilt mock data endpoints    |
|  ConfigHandler    -> Field types & configs           |
|-----------------------------------------------------|
|        MongoDB          Redis          Gin API       |
+-----------------------------------------------------+


⚙️ Setup Instructions

1️⃣ Clone the Repository
git clone https://github.com/saifwork/mock-service.git
cd mock-service

2️⃣ Install Dependencies
go mod tidy

3️⃣ Add .env File
Create a .env in the root directory:

APP_NAME=Mock Service
APP_PORT=8080
APP_ENV=development

# MongoDB
MONGO_URI=mongodb+srv://<username>:<password>@cluster.mongodb.net/?retryWrites=true&w=majority
MONGO_DB_NAME=mock_service

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=supersecret

# Gmail
GMAIL_USER=your_email@gmail.com
GMAIL_PASSKEY=your_app_password

4️⃣ Run the App
go run main.go


You should see:

[MONGO] Connected successfully
[GIN] Listening on port 8080...


--------------------------------------------------------------

🧩 API Endpoints

# 🧑‍💻 Auth Routes
Method	Endpoint	Description

POST	/api/auth/signup	Register new user
POST	/api/auth/login	Login
GET	/api/auth/verify-email?token=	Verify email
POST	/api/auth/forgot-password	Send reset email
POST	/api/auth/reset-forgot-password	Reset forgotten password
POST	/api/auth/reset-password	Change password (JWT protected)

# 📦 Project Routes
Method	Endpoint	Description
POST	/api/projects	Create project

GET	/api/projects	List user projects
GET	/api/projects/:pid	Get project by ID
PUT	/api/projects/:pid	Update project
DELETE	/api/projects/:pid	Delete project

# 📚 Collection Routes
Method	Endpoint	Description

POST	/api/projects/:pid/collections	Create new collection
GET	/api/projects/:pid/collections	Get all collections
GET	/api/collections/:cid	Get collection by ID
PUT	/api/collections/:cid	Update collection
DELETE	/api/collections/:cid	Delete collection

# 🗂️ Record Routes
Method	Endpoint	Description

POST	/api/collections/:cid/records	Create record
GET	/api/collections/:cid/records	List records
GET	/api/records/:rid	Get record by ID
PUT	/api/records/:rid	Update record
DELETE	/api/records/:rid	Delete record

# ⚙️ Config Routes
Method	Endpoint	Description

GET	/config/types	Supported field types
GET	/config/preset/users	Prebuilt mock users
GET	/config/preset/products	Prebuilt mock products
GET	/config/preset/posts	Prebuilt mock posts
GET	/config/preset/comments	Prebuilt mock comments
GET	/config/preset/carts	Prebuilt mock carts
GET	/config/preset/todos	Prebuilt mock todos

# 🧪 Health
Method	Endpoint	Description

GET	/health	Check service status


🧑‍🎨 Author ---

Md Saif
Software & Game Developer | Golang + Flutter Enthusiast
📍 GitHub

💬 Passionate about developer tools & backend systems