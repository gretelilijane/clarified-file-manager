# Quick set-up with docker

```sh
docker compose up --build
```

# Alternative development set-up

## Requirements
- Node.js (for development only)
- Go 1.23.3
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
