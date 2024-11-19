# Quick set-up with docker

```sh
docker compose up --build
```

# Alternative development set-up

## Requirements
- Node.js
- Go 1.16
- PostgreSQL

## Environment

```sh
cp .env.example .env
```
Set up postgres database and fill in the environment variables in `.env`

## Install dependencies

```sh   
npm install
``` 

## Run 

```sh
npm run dev
```
