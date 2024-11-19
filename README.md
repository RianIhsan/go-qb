

#  go-qb is Query Builder for Go

## Description

**Query Builder** is a simple tool for building safe SQL queries in the Go programming language. Designed to prevent SQL injection risks, this Query Builder supports two popular database types: PostgreSQL and MySQL. By using appropriate placeholders for each database type, this tool allows users to easily and securely construct complex SQL queries.

## Key Features

- **Support for Multiple Databases**: Supports queries for both PostgreSQL and MySQL.
- **Dynamic Query Building**: Makes it easy to create SQL queries with dynamic conditions using the `Param` struct.
- **SQL Injection Prevention**: Uses parameterized queries to enhance security and reduce the risk of SQL injection.
- **Flexibility**: Offers functions for `SELECT`, `UPDATE`, and `DELETE`, along with the ability to easily add additional queries.
- **Dynamic Placeholders**: Automatically generates placeholders based on the type of database being used.


# Install
```bash
go get github.com/RianIhsan/go-qb
```

## How to use ?

- Example:
    <br> - user_repository.go

    ```golang
    func (db *ReadWrite) GetUsers(ctx context.Context, req md.ReqGetUsers) (resp []md.RespGetUsers, err error) {
        qb := querybuilder.New(querybuilder.DBPostgres, DropYourQueryHere)
        
        if req.Age != "" {
            qb.AddQuery(" AND age = ? ", req.Age)
        }
        qb.AddString(`ORDER BY name ASC`)

        // ......
    }
  ```