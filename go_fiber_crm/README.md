# Folder structure

```
.
├── database
│   └── database.go
├── go.mod
├── go.sum
├── lead
│   └── lead.go
├── main.go
└── README.md
```

1. database - this can be inside config folder. this handles the opening and creating the database connections
2. lead - this should be seperated between controllers and  models folder. This handles the data and controllers.

# Content
This project is a simple customer relationship management system.
This project uses **fiber**, **sqlite** and **gorm**.
This project explores a basic usage of **fiber** and **sqlite**
