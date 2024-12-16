
# That Pet Place Backend

This repository contains the backend for **That Pet Place**, a platform that connects pet owners with shops, clinics, and service providers. The backend is built using Go and offers RESTful APIs to support platform functionality.

## Overview

That Pet Place Backend serves as the core for managing:

1. **Pet Owners**: Users can register, add their pets, search for services, and book appointments or purchase products.
2. **Service Providers**: Shops, clinics, or service providers can manage products, appointments, and respond to user requests.

---

## Key Features

### For Pet Owners:
- **Account Management**: Register and log in securely.
- **Pet Management**: Add, update, or remove pet profiles.
- **Service Search**: Find nearby shops, clinics, or service providers based on location.
- **Appointment Booking**: Book appointments for services or purchase products.

### For Service Providers:
- **Dashboard Management**: View and manage appointments and bookings.
- **Product and Service Management**: Add and update services or products.
- **Appointment Handling**: Accept or cancel appointment requests.

---

## System Architecture

The backend uses a modular architecture, separating responsibilities into distinct components for better maintainability and scalability.

### High-Level Flow
1. **User Interaction**: Users interact via RESTful APIs.
2. **Authentication**: JWT-based authentication ensures secure access.
3. **Database Operations**: PostgreSQL and MongoDB handles data for pets, users, services, and appointments.
4. **Business Logic**: Encapsulated in the service layer for clear separation of concerns.

---

## Technologies Used

- **Language**: Go 
- **Database**: PostgreSQL, MongoDB
- **Authentication**: JWT (JSON Web Tokens)
- **Environment Management**: dotenv for configuration
- **Dependency Management**: Go Modules

---

## Database Schema

### Key Entities

- **Users**: Stores user information (ID, name, email, password).
- **Pets**: Represents pets associated with users (ID, name, type, age, owner ID).
- **Services**: Details services provided by shops or clinics (ID, name, location, provider ID).
- **Appointments**: Tracks booking details (ID, user ID, service ID, status, date).

### Relationships
- **Users ↔ Pets**: One-to-many (a user can own multiple pets).
- **Users ↔ Appointments**: One-to-many (a user can book multiple appointments).
- **Services ↔ Appointments**: One-to-many (a service can have multiple bookings).

---



## Project Structure

```plaintext
that-pet-place-backend-go/
├── cmd/                # Main application entry point
├── config/             # Configuration files and utilities
├── services/           # Business logic for each service (users, pets, etc.)
├── types/              # Types and interfaces
├── utils/              # Utility functions (e.g., helpers, formatters)
├── Makefile            # Build and run commands
└── go.mod              # Go module file
```



# Frontend: 
https://github.com/Farmaan-Malik/ThatPetPlace

