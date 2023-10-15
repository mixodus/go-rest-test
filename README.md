<!-- GETTING STARTED -->
### Built With
Go: go version go1.21.3 darwin/arm64
* <a href="https://go.dev/">GO<a>
* <a href="https://gin-gonic.com/">GIN<a>
* Uses PostGreSQL
  
## Getting Started

Clone the project

```
git clone https://github.com/mixodus/go-rest-test.git
```

There are few things that need to be setup before running `go run main.go` command at the root project `cd go-rest-test`.

### Prerequisites

* Setup Database
  
  - Create new database, and feel free to name the database.
  - Restore or Import  `go_test.sql` /  `go_test.backup` to your new database.

  We need to import DB for the Banks data. I only give two row of Banks, if you need more than that you can manually inject new row of Banks.
  
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



Cheers!!🥂✨




### Explanation

* `Login` uses redis to store session, if player re-login then player need to use new token else unauthorized.
* Same as for `Logout`, uses redis. If player logout then no token in redis for check, so in any ways players always will be unautorized untill player do login again.
* For `Wallet` I create `transaction` table to count and refresh player's balance. `GET` API for `Player's Wallet` will always re-sum or count between `transaction.transaction_type` (DEBIT/CREDIT) and restore it in `players.balance` column
* Players need to have Bank information in order to make `Top Up` or transactions.
* Players only can have ONE bank information. If want to update or make a new one, player need to delete previous bank by hit/consume `DELETE Player Bank` API.
* `transaction` has statuses, re-sum or count between `transaction.transaction_type` (DEBIT/CREDIT) only applies on `success` statuses. I've make public APIs on pretend to be `ADMIN` for changing statuses to `success`. APIs are `Set Debit Success` && `Set Credit Success`
* Public APIs on pretend to be `ADMIN` is for getting player list too. `GET Players`
* `Top Up` transaction request need file to be uploaded, file path will stored at `players_banks.file_name` table. You can try to access file path manually from DBMS or pgAdmin4 then copy it and access the file source via ```http://localhost:8080/api/image?path=<file_path>``` !! .uploads folder is `gitignored` !! folder will automatically created if not exist.

### Postman Collection Preview
<img width="288" alt="Screenshot 2023-10-15 at 18 20 44" src="https://github.com/mixodus/go-rest-test/assets/58242458/f8beaa15-4d0d-4040-a4f4-d3796399cdc8">
