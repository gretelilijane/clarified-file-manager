# About
# Set-up
## Docker
```sh
sudo docker-compose up --build
```
## Dev
### Install go modules
```sh
go mod tidy
```
### Database
1. Install PostgreSQL
2. Set up password for user `postgres`
3. Allow password authentication in `pg_hba.conf`
4. Run 
```sh
.\setup_db.sh
```
to create the database and tables
### Run

```sh   
npm install
``` 
```sh
npm run start-dev
```