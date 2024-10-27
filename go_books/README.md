# Folder structure
```
.
├── cmd
│   └── main
│       └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
└── pkg
    ├── config
    │   └── app.go
    ├── controllers
    │   └── book-controller.go
    ├── models
    │   └── book.go
    ├── routes
    │   └── bookstore-routes.go
    └── utils
        └── utils.go
```

1. config - the connection with the database is established here.
2. conrollers - route controllers
3. models - communicate with db and handle data
4. routes - routes of the app
5. utils - utility functions, mostly pure functions

# Content
This project has a proper folder structure. 
It is a very simple book management api with **MySQL**.
It uses **Gorilla Mux** and **Gorm** ORM.
