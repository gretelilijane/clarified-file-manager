# About
# Set-up
## Docker
```sh
sudo docker-compose up --build
```
## Dev
### Nodemon
```sh
nodemon --exec "go run main.go" --signal SIGTERM -e go,env,html
```