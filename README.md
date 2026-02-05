# 📝 Todo-Go: Clean Architecture Task Manager

<details>
<summary><h2><strong>Table of Contents</strong></h2></summary>

- [About the Project](#about-the-project)
- [Tech Stack 🚀](#tech-stack-)
- [Key Features](#key-features)
- [Project Architecture](#project-architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation & Docker Setup](#installation--docker-setup)
- [Usage](#usage)
- [Testing](#testing)
- [Deployment & Security](#deployment--security)

</details>

## About the Project

**Todo-Go** is a high-performance, full-stack task management application designed with a focus on **Clean Architecture** and scalability. It features a robust Go backend and a responsive React frontend, providing users with a seamless experience for organizing tasks via custom lists, labels, and drag-and-drop functionality.

## Tech Stack 🚀

### Backend (The Root)
* **Language:** Go (Golang) 1.25
* **Routing:** Chi Router v5
* **Database:** PostgreSQL with `sqlx` for struct scanning
* **Authentication:** JWT (JSON Web Tokens) with `jwtauth`
* **Migrations:** `golang-migrate` for versioned schema control

### Frontend (`/client`)
* **Framework:** React 19 + Vite
* **Styling:** Tailwind CSS + DaisyUI
* **State Management:** React Context API (Auth & Lists)
* **Interactions:** `@dnd-kit` for drag-and-drop task reordering

### Infrastructure
* **Proxy:** NGINX (Production reverse proxy)
* **Containerization:** Docker & Docker Compose

---

## Key Features

* **Secure Authentication:** User registration and login using JWT and `bcrypt` password hashing.
* **Custom Todo Lists:** Create multiple lists with specific titles, colors, and category labels.
* **Task Lifecycle:** Add, edit, toggle, and delete individual todos within lists.
* **Soft Delete (Bin):** Move lists to a trash bin and restore them later or delete them permanently.
* **Global Search:** Filter lists and items in real-time via the sidebar search.
* **N+1 Optimization:** Backend DTOs support fetching lists with nested items in a single request.

---

## This is the home page:

<img width="1485" height="789" alt="Image" src="https://github.com/user-attachments/assets/293f75c4-40e7-4e5b-9fb1-f41fe90f2c05" />

## Project Architecture

The backend follows **Clean Architecture** to decouple business logic from infrastructure:

1.  **Domain Layer:** Core entities (`User`, `Todo`, `TodoList`) and business errors.
2.  **Service Layer:** Business rules, validation, and orchestration.
3.  **DAL (Data Access Layer):** SQL execution using template-based query building.
4.  **Delivery Layer:** HTTP handlers and JWT middleware.

---

## Getting Started

### Prerequisites
* [Docker Desktop](https://www.docker.com/) installed and running.

### Installation & Docker Setup

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/macesz/todo-go](https://github.com/macesz/todo-go)
    cd todo-go
    ```

2.  **Configure Environment Variables:**
    Create a `.env` file in the root directory:
    ```env
    DB_USER=postgres
    DB_PASS=password
    DB_NAME=go_todo
    JWT_SECRET=your_super_secret_key
    SERVER_PORT=3000
    ```

3.  **Launch with Docker Compose:**
    ```bash
    docker compose up --build
    ```
    This command builds the Go binary, runs database migrations, compiles the React frontend, and starts the NGINX reverse proxy.

4.  **Access the Application:**
    * **App:** [http://localhost:8080](http://localhost:8080)
    * **API Health:** [http://localhost:8080/api/health](http://localhost:8080/api/health)

---

## Usage

1.  **Auth:** Register a new account to receive a secure JWT.
2.  **Lists:** Use the "Take a note..." bar to create new categorized lists.
3.  **Tasks:** Click a task to edit its title or drag it to reorder within its priority group.
4.  **Management:** Rename or delete labels globally via the "Edit Labels" modal in the sidebar.

---

## Testing

The project emphasizes test-driven development with high coverage:

* **Unit Tests:** `go test -short ./...`
* **Integration Tests:** `go test -v ./tests` (Requires Docker for Testcontainers)
* **Frontend Linting:** `cd client && npm run lint`

---

## Deployment & Security

* **NGINX Reverse Proxy:** Serves frontend assets and routes `/api` calls, hiding the backend port for security.
* **Password Security:** Plaintext passwords never enter the database; they are hashed with a cost factor of 10.
* **Case Sensitivity:** Handlers and build stages are optimized for Linux-based container environments.
