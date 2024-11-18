# SQlite
```sh
go get github.com/mattn/go-sqlite3
```

## Gorilla Mux
```sh
go get -u github.com/gorilla/mux
```

## Gorilla sessions
```sh
go get -u github.com/gorilla/sessions
```

## Godotenv
```sh
go install github.com/joho/godotenv/cmd/godotenv@latest
```

## Postgres
```sh
go get -u github.com/lib/pq
```

## Argon2
```sh
go get -u golang.org/x/crypto/argon2
```

# Run
```sh
go run main.go
```

# Set-up potgreSQL

## Install
```sh
sudo apt update
sudo apt install postgresql postgresql-contrib
```

## set-up postgres user
```sh
sudo -i u postgres
psql
``` 

## Connect with username to Postgres
```sh
psql -U postgres
```

### Add user password
```sql
ALTER USER postgres PASSWORD 'password';
```

### Modify pg_hba.conf to allow password authentication
```sh
sudo nano /etc/postgresql/<version>/main/pg_hba.conf
```

Change the following line
```sh
local   all             postgres                                peer
```
to 

```sh
local   all             postgres                                md5
```
Restart postgres
```sh
sudo systemctl restart postgresql
```

Test connecting to postgres with password
```sh
psql -U postgres -W
```

### Create a new user, db and give user access to db
```sql
CREATE DATABASE new_db;
CREATE USER username WITH PASSWORD 'securepassword';
GRANT ALL PRIVILEGES ON DATABASE new_db TO username;

\c new_db

GRANT USAGE ON SCHEMA public TO username;
GRANT CREATE ON SCHEMA public TO username;
```

### Connect to db
```sh
psql -U username -d new_db -W
```

### add tables to the db
```sh
psql -U username -d clarified_file_manager_db -W -f db/init.sql
```

## Dev

### Nodemon
```sh
nodemon --exec "go run main.go" --signal SIGTERM -e go,env,html
`
## Used literature

### PostgreSQL

#### Creating DB
- https://www.calhoun.io/creating-postgresql-databases-and-tables-with-raw-sql/

#### Connecting to DB with GO
- https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/


# Login with go and htmx:
https://github.com/guillemaru/authentication-playground/tree/master

# Password hashing and salting
https://snyk.io/blog/secure-password-hashing-in-go/