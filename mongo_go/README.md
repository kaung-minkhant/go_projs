# Folder structure

```
.
├── controllers
│   └── users.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── models
│   └── users.go
└── README.md
```

1. controller - route controllers
2. models - data models

# Content

This project builds a simple user management api.
This uses **httprouter** and official mongodb driver.

This project uses **interfaces** to model controllers.
This project does the data handling in controller, which should not be done.
