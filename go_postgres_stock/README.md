# Folder structure

```
.
├── controllers
│   └── handlers.go
├── database
│   └── postgres.go
└── router
    └── router.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── models
│   └── models.go
├── README.md
```
1. controllers - route controllers
2. database - handle database connections
3. router - app routes
4. models - data models

# Content

This project usess **Postgresql** db and **gorilla mux**.

# Notes
Be sure to create db and table in postgres before running the application
