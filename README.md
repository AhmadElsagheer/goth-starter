# 🛡️ Goth-Starter

**The BetterAuth.TS + Go Backend Starter**

A reusable authentication boilerplate for solo-founders and bootstrappers building in Go and need a solid auth solution like BetterAuth.

> [!NOTE]
> This is a **boilerplate**, not a package. You are meant to clone it, eject it, and modify it. You own the code.

## 🏗️ Architecture

```mermaid
graph TD
    Client[Client]
    Auth["BetterAuth.TS<br/>(Auth Server)"]
    Resource["Go Backend<br/>(Resource Server)"]

    subgraph DB [Database]
        direction LR
        AuthTables["Auth Tables<br/>(users, sessions, accounts, verifications)"]
        BusinessTables["Business Tables<br/>(business_table_1 ... business_table_n)"]
    end

    %% Force Layout: Auth Left of Resource
    Auth ~~~ Resource

    %% Auth Flow
    Client -- 1. Authenticate --> Auth
    Auth -- 2. Returns accessToken (Bearer) --> Client
    
    %% Request Flow
    Client -- 3. Business Request + accessToken --> Resource
    
    %% Data Access
    Auth -- Read/Write --> AuthTables
    Resource -- Read-Only --> AuthTables
    Resource -- Read/Write --> BusinessTables
```

1.  **Auth Server (BetterAuth.TS)**: A standalone Bun instance running [BetterAuth](https://www.better-auth.com). It handles all writes to the auth tables and manages OAuth/Email/Phone flows.
2.  **Resource Server (Go)**: Your main backend. It treats the auth tables as **Read-Only**. It uses the `go-kit` middleware to validate session tokens directly against the database with zero network overhead.
3.  **Database (Postgres)**: A shared "Unified Database" where both the Auth Server and Go Backend connect. No data syncing required.

## 📁 Repository Structure

- **`/auth-server`**: The heart of authentication. Contains the BetterAuth instance, configuration, and auth db schema.
- **`/go-kit`**: Drop-in Go code. Includes Echo middlewares for auth and RBAC permissions and db access.
- **`/react-kit`**: Frontend starter kit. Includes Tailwind-ready login/signup components and hooks for session management.

## 🚀 Quick Start

⚡ [Run the quick start example](./examples/gother-example)

### Before you start
- Required tools: go, psql, [sqddl](https://github.com/bokwoon95/sqddl)
- Go dependencies in the `go-kit`: zap, echo, pgx, sq, fx (you can change or remove them after generation)

### Install
1. Set variables
```bash
export USERNAME=AhmadElsagheer
export PROJECT_NAME=gother-example
export GO_BACKEND_PATH=./examples/gother-example/backend
export AUTH_SERVER_PATH=./examples/gother-example/auth-server
export MODULES_PATH=server/modules
export SCHEMA_PATH=server/schema/postgres/types
export MIGRATIONS_PATH=server/schema/postgres/migrations
export DB_URI="postgres://${PROJECT_NAME}:${PROJECT_NAME}@localhost/${PROJECT_NAME}?sslmode=disable"
```

2. Create Project
```bash
./create-go-project.sh $USERNAME/$PROJECT_NAME $GO_BACKEND_PATH
```

3. Run database
```bash
docker compose -f $GO_BACKEND_PATH/docker-compose.yaml up -d
```

4. Setup auth server
```bash
export BETTER_AUTH_SECRET=$(openssl rand -base64 32)
cp -r auth-server $AUTH_SERVER_PATH && (cd $AUTH_SERVER_PATH && bun install)
(cd $AUTH_SERVER_PATH && cp .env.example .env && \
sed -i '' "s|^BETTER_AUTH_SECRET=.*|BETTER_AUTH_SECRET=$BETTER_AUTH_SECRET|" .env && \
sed -i '' "s|^DATABASE_URL=.*|DATABASE_URL=$DB_URI|" .env)
```

5. Generate the migration file (NOTE: database must be running. We also move the migration file to the backend project to be co-located with other business table migrations)
```bash
mkdir -p $GO_BACKEND_PATH/${MIGRATIONS_PATH}
export MIGRATION_FILE_ABS_PATH=$(cd "$GO_BACKEND_PATH/${MIGRATIONS_PATH}" && pwd)/001_create_auth_tables.up.sql
(cd $AUTH_SERVER_PATH && bun x @better-auth/cli@latest generate --output $MIGRATION_FILE_ABS_PATH -y)
```

6. Migrate the database
```bash
psql -h localhost -U $PROJECT_NAME -d $PROJECT_NAME -f $MIGRATION_FILE_ABS_PATH
```

7. Install go kit
```bash
./install-go-kit.sh \
  --project $GO_BACKEND_PATH \
  --modules-path $MODULES_PATH \
  --db-uri $DB_URI \
  --schema-path $SCHEMA_PATH
```

8. Implement `$GO_BACKEND_PATH/cmd/server/main.go`. You can use the example in `$GO_BACKEND_PATH/examples/gother-example/backend/cmd/server/main.go`.

9. Set env variables in `$AUTH_SERVER_PATH/.env`
- `AUTH_METHODS`: Comma-separated list (e.g., `email,phone,google`).
- `SMTP_*`: For email verification.
- `TWILIO_*`: For phone number OTPs.
- `GOOGLE_*`: For Google OAuth.

### Run
```bash
# Run the auth server
cd $AUTH_SERVER_PATH && bun run dev

# Run the go backend
cd $GO_BACKEND_PATH && go run ./cmd/server
```
### Test
Use [bruno](https://www.usebruno.com/) to test the API. You can find the bruno requests in [`.bruno` workspace](./.bruno).

Auth Server requests:
1. `Email - Signup`: returns 200 OK. The verification code is sent to the email address.
2. `Email - Verify OTP`: edit the code in request body and send the request. It returns 200 OK and the user is verified.
3. `Email - Signin`: returns 200 OK and session token in response
4. `Get Token`: sends session token to auth server and returns jwt access token (this is what the backend accepts).

Backed Server requests:
5. `Ping`: returns 200 OK and user info.
6. `Admin Ping`: returns 401 (insufficient permissions) since our user is not admin. You can add the `users:read` permission to the customer role in `server/modules/auth/rbac.go` to allow the customer to access this endpoint.



## 🧠 Workflows & Features

### Authentication Flows
- **Email/Password**: Full flow supported: Sign up -> Verify Email (OTP) -> Sign in.
- **Phone Auth**: Sign up -> Sign in via SMS OTP.
- **Social Auth**: Google, GitHub, etc. (Configurable in `auth.ts`).

### Clients
- **Mailer**: Built-in support for simple SMTP email sending.
- **SMS**: Built-in Twilio integration for phone OTPs.

### Credentials
- SMTP: [use google workspace email (or gmail) + app password](https://support.google.com/accounts/answer/185833?hl=en) or change the mailer implementation to use [oauth2 and gmail api](https://developers.google.com/workspace/gmail/api/quickstart/nodejs)
- [Twilio](https://www.twilio.com/en-us/blog/send-sms-twilio-shell-script-curl)
- [Google](https://www.better-auth.com/docs/authentication/google#get-your-google-credentials)


## 📝 Notes
- **Seeding**: You can seed initial users with migrations in backend. You can find example data in `examples/seed-users.json`.
- **Temp Emails**: Email is required for all users. BetterAuth assigns temporary emails to users who signup with phone. You can change the format in `auth-server/src/lib/auth.ts`.
- **Templates**: Content for Email and SMS can be customized in `auth-server/src/lib/templates.ts`.
