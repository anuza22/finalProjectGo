# chatAppGo
# 💬 chatAppGo

A lightweight, modular chat backend written in Go.  
Built with clean architectural principles for scalability, clarity, and maintainability.

---

## 📌 Overview

**chatAppGo** is a backend service that implements core chat functionality such as user authorization and request handling. The application is structured with a layered and modular architecture that ensures maintainability, testability, and clear separation of concerns.

---

## ✨ Features

- 🔐 JWT-based user authentication
- 🧱 Clean, layered architecture
- 🌐 HTTP handling using Go's `net/http`
- 📦 Modular package structure
- 🧩 Ready for extension: REST APIs, database, WebSocket, etc.

---

## 📁 Project Structure

chatAppGo/ ├── cmd/ # Application entry point │ └── main.go # Startup logic │ ├── Authorization/ # Authentication and token logic │ ├── auth.go # HTTP handler for login/auth │ └── token.go # JWT generation and validation │ ├── internal/ # Core business logic │ └── user/ # Example domain logic (user module) │ ├── pkg/ # Shared utilities │ └── config/ # Configuration loading │ ├── test/ # Test files (optional) │ ├── go.mod # Go module file ├── go.sum # Dependency checksums └── README.md # Project documentation

> ✅ Follows a layered architecture:  
> `Handler` ➝ `Service` ➝ `Domain` ➝ `Infrastructure`

---

## ⚙️ Getting Started

### ✅ Requirements

- Go 1.21 or higher
- Git

### 📥 Installation

Clone the repository:

```bash
git clone https://github.com/anuza22/finalProjectGo.git
cd finalProjectGo
go mod tidy
▶️ Run the Application
go run ./cmd/main.go
🧪 Testing

To run all unit tests:
go test ./...
⚠️ Tests will be added as development progresses
🛠 Technologies

Component	Description
Go	Backend language
net/http	HTTP server
JWT	Token-based auth
Modular Arch	Layered codebase
🚀 Possible Extensions

 WebSocket for real-time chat
 Database integration (PostgreSQL or SQLite)
 RESTful API endpoints
 Middleware for logging, CORS, tracing
 Docker support for containerization
 Full unit and integration test coverage
👤 Authors

Developed by Anuza, Alikhan, Ayazhan
Final project @ KBTU University — Go Development Course
📄 License

This project is licensed under the MIT License.
Feel free to use, modify, and distribute it.
