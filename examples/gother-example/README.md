# gother-example


## Install
```sh
cd auth-server
cp .env.example .env
# REQUIRED: fill the following variables in .env (check repo README.md for how to get them)
# - BETTER_AUTH_SECRET
# - SMTP_USER
# - SMTP_PASSWORD
bun install
cd ..
```

## Run servers
```sh
# T1: Run auth server
cd auth-server && bun run dev

# T2: Run go backend
cd backend && go run ./cmd/server
```
