<!-- GETTING STARTED -->
### Built With
Go: go version go1.21.3 darwin/arm64
* <a href="https://go.dev/">GO<a>
* <a href="https://gin-gonic.com/">GIN<a>
  
## Getting Started

clone project

```
[git clone ](https://github.com/mixodus/go-rest-test.git)
```

There are few things that need to be setup before running `go run main.go` command at the root project.

### Prerequisites

* Setup Database
  1. We use PostGreSQL. Create new database, and feel free to name the database.
  2. Restore or Import  `go_test.sql` /  `go_test.backup` to your new database.

  We need to import DB for the Banks data. I only give two row of Banks, if you need more than it you can manually inject new row of Banks.
  
* Setup `.env` file
  Change it to your database information.
  ```
  DATABASE_HOST = "localhost"
  DATABASE_PORT = "5432"
  DATABASE_USER = "postgres"
  DATABASE_PASSWORD = "ivanyunus"
  DATABASE_NAME = "go_test"
  ```

* Import Postman Collection
  Import `GO.postman_collection.json` to your desired software to test the APIs.

