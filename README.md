# ITK_Academy-Test
Solution for the test task of the ITK Academy company

This project consists of two components:

backend — the Go-based API server
postgres — a PostgreSQL database

Make sure you have the following installed:
- Docker
- Docker Compose

## 📁 Project Structure
```
├── cmd/
│   └── config.env/
├── config/
├── internal/
│   ├── dto/
|   ├── handlers
│   ├── models
│   ├── repository
│   └── services
├── tests/
│   ├── handlers/
│   ├── repositories/
│   └── services/
├── gitignore
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## ⚙️ How to Run
1. Clone the repository:
```
git clone https://github.com/Mimist-Illusionard/ITK_Academy-Test.git
```
2. Start the project:
```
docker-compose up --build
```
This command will build and run all three containers:

go-blog-backend will be accessible at: http://localhost:9090

go-blog-postgres will run PostgreSQL on port 5432

3. Stopping the project:
```
docker-compose down
```

## Postman Collection
[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://app.getpostman.com/run-collection/44290956-c1dc944a-58f2-48f5-9f1d-5d571b87e57a?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D44290956-c1dc944a-58f2-48f5-9f1d-5d571b87e57a%26entityType%3Dcollection%26workspaceId%3D373e624b-b49c-43d7-9b00-b7dbb0ed6baa)