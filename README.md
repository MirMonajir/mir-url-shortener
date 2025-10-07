# URL Shortener Service

A simple URL shortener REST API built in **Golang** using the **Gin** framework.  
The service shortens URLs, resolves short URLs to their original form, and provides domain usage metrics.

---

## Features

- Shorten URLs with **unique 6-character random codes**.
- Reuse existing short URL if the original URL was shortened before.
- Redirect short URLs to the original URLs.
- Track and return **top 3 most shortened domain names**.
- Store all data **in-memory**.
- RESTful API design.
- Unit tested with clear separation of concerns (DDD & SOLID principles).
- Dockerized for easy deployment.

---

## Tech Stack

- Go 1.25+
- Gin Web Framework
- Docker & Docker Compose
- Go Modules

---

## Getting Started

### Prerequisites

- Go (1.25 or higher) installed: https://golang.org/dl/
- Docker installed (optional for containerized deployment)

---

### Running Locally

1. Clone the repo:

   ```bash
   git clone https://github.com/<your-username>/mir-url-shortener.git
   cd mir-url-shortener

2. Run the server:
   ```bash
   go run ./cmd/main.go

3. The server will start on http://localhost:8080


### Running with Docker

1. Pull the Docker image from Docker Hub:

   ```bash 
     docker pull mirmonajir/mir-url-shortener:latest


2. Run the container:

   ```bash 
      docker run -p 8080:8080 mirmonajir/mir-url-shortener:latest


Access the service at:

http://localhost:8080
