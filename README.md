# That Pet Place Backend

This repository contains the backend for **That Pet Place**, a platform that connects pet owners with shops, clinics, and service providers. It is written in **Go**, follows **Clean Architecture principles**, and provides **secure RESTful APIs** to power the platform.

## 🚀 Overview

That Pet Place Backend is the core system that enables:

* **Pet Owners**: Register, manage pets, search for services, and book appointments or buy products.
* **Service Providers**: Manage services, products, and appointments through a secure provider dashboard.

---

## ✨ Key Features

### 🔹 For Pet Owners

* **Account Management**: Register, log in, and authenticate with JWT.
* **Pet Management**: Add, update, and remove pet profiles.
* **Service Discovery**: Search nearby shops, clinics, or providers.
* **Appointments & Purchases**: Book services or buy products.

### 🔹 For Service Providers

* **Dashboard**: Manage appointments and bookings.
* **Service & Product Management**: Add and update offerings.
* **Appointment Handling**: Accept, decline, or manage requests.

---

## 🏗️ System Architecture

The project uses **Clean Architecture** with clear separation of concerns:

1. **Routing Layer** – Gorilla Mux for request routing.
2. **Authentication** – JWT-based authentication & middleware.
3. **Business Logic Layer** – Encapsulated in services (use cases).
4. **Persistence Layer** – PostgreSQL (structured data) + MongoDB (unstructured/flexible data).

---

## 🛠️ Tech Stack

* **Language**: Go
* **Framework**: Gorilla Mux
* **Architecture**: Clean Architecture
* **Databases**: PostgreSQL + MongoDB
* **Auth**: JWT (JSON Web Tokens)
* **Configuration**: dotenv
* **Dependency Management**: Go Modules

---

## 🗄️ Database Schema (High-Level)

### Entities

* **Users**: Profile and credentials (Postgres).
* **Pets**: Linked to owners (Postgres).
* **Services**: Shops/clinics with location & offerings (MongoDB for flexible service data).
* **Appointments**: Bookings linking users & providers (Postgres).

### Relationships

* `User → Pets` (1\:N)
* `User → Appointments` (1\:N)
* `Service → Appointments` (1\:N)

---

## 📂 Project Structure

```plaintext
that-pet-place-backend-go/
├── cmd/            # Application entry points (API setup)
├── config/         # Configuration & environment handling
├── store/          # Data access layer for Postgres & MongoDB
├── services/       # Business logic & use cases
├── types/          # Shared types and interfaces
├── utils/          # Utilities (JWT, helpers, etc.)
├── go.mod          # Go module definition
└── Makefile        # Build/run commands
```

---

## ▶️ Getting Started

### Prerequisites

* Go 1.23+
* PostgreSQL & MongoDB running locally or via Docker

### Setup

```bash
# Clone repo
git clone https://github.com/lunatictiol/that-pet-place-backend-go.git
cd that-pet-place-backend-go

# Install dependencies
go mod tidy

# Run the server
make run
```

Server will start on `http://localhost:8080`.

---

## 🌐 Frontend

The React frontend for this project is available here:
👉 [That Pet Place Frontend](https://github.com/Farmaan-Malik/ThatPetPlace)


