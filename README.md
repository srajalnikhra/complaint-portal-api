# Complaint Portal API

A RESTful Complaint Portal API built using Go's standard `net/http` package.

This project demonstrates user registration, authentication, complaint management, request validation, JSON handling, concurrency using mutexes, and REST API design without using any third-party frameworks.

---

## Features

- User Registration
- User Login
- Generate Secret Code
- Submit Complaint
- View User Complaints
- View All Complaints (Admin)
- View Complaint by ID
- Resolve Complaint
- In-memory Data Storage
- Mutex for Concurrent Requests
- JSON Request & Response Handling

---

## Tech Stack

- Go
- net/http
- encoding/json
- sync
- math/rand

---

## Project Structure

```
complaint-portal-api/
│
├── go.mod
├── main.go
└── README.md
```

---

## API Endpoints

| Method | Endpoint | Description |
|---------|----------|-------------|
| POST | /register | Register User |
| POST | /login | User Login |
| POST | /submitComplaint | Submit Complaint |
| GET | /getAllComplaintsForUser | View User Complaints |
| GET | /getAllComplaintsForAdmin | View All Complaints |
| GET | /viewComplaint | View Complaint by ID |
| PUT | /resolveComplaint | Resolve Complaint |

---

## Run Locally

```bash
git clone https://github.com/srajalnikhra/complaint-portal-api.git
```

```
cd complaint-portal-api
```

```
go run main.go
```

Server starts on

```
http://localhost:8080
```

---

## Learning Objectives

- REST API Development
- User Authentication
- CRUD Operations
- JSON Encoding & Decoding
- Mutex Synchronization
- Request Validation
- Go Standard Library