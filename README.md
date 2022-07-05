
# Inventory-Management-System

A web-based Inventory Management System application written in GoLang.


## Installation

Clone the repository

```bash
git clone https://github.com/mohd4hassan/DIT-Inventory-Management-System.git
```
Initialise ```go mod``` at the root directory of the project 
```bash
go mod init IMS
```
Install missing packages and libraries
```bash
go mod download
```
then run 
```bash
go mod tidy
```
## Database configuration
- Download and install PostgreSQL DBMS
- Create a new Database called `IMS`
- Username: `user` and Password: `admin`

**Alternatively**: You can modify the config files to your preference.

`config.go`
```go
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "admin",
		Name:     "IMS",
	}
}
```
`.config`
```go
    . . .
    "database": {
        "host": "localhost",
        "port": 5432,
        "user": "user",
        "password": "admin",
        "name": "IMS"
    }
```
## Run the program
To run the program, install **Fresh**. GoLang webserver builder and refresh (restart) tool.
```go
go get -u github.com/pilu/fresh
```
#### Usage
Make sure the you are at the root directory of the project, then run:
```go
fresh
```
Fresh will watch for file events, and every time you create/modify/delete a file it will build and restart the application. If `go build` returns an error, it will log it in the tmp folder.

Open the application on your browser:
```go
http://localhost:3030/
```

## Assistance
For any assistance on this project, email to [Mohamed Hassan](mailto:mohdalhassan2012@gmail.com "Mohamed Hassan").