# Go Expense REST API

A simple API to create, list, and update expenses. Go runs the API locally. Docker Compose runs PostgreSQL for the database setup.

## 1. Install Go and Docker

Install Go 1.27.3 or newer from [Go downloads](https://go.dev/dl/). Install Docker with the Compose plugin if you plan to run PostgreSQL.

```bash
go version
docker --version
docker compose version
```

## 2. Download the project

Replace the example URL with this project's Git URL.

```bash
git clone 'https://example.com/your-repo.git' go-lang-restapi
cd go-lang-restapi
```

If you downloaded the parent learning repository, open its `go-lang-restapi` folder instead. Run the remaining setup commands from the folder containing `go.mod` and `compose.yaml`.

## 3. Start PostgreSQL and apply the migration

The commands assume your Compose service is named `postgres`, with user `expense_user` and database `golang_expense_tracker`.

```bash
docker compose up -d postgres
docker compose ps postgres
docker compose exec postgres pg_isready -U expense_user -d golang_expense_tracker
```

Once PostgreSQL reports `accepting connections`, apply the migration and verify the table.

```bash
docker compose exec -T postgres psql -U expense_user -d golang_expense_tracker \
  -v ON_ERROR_STOP=1 < migrations/001_create_expenses.sql

docker compose exec postgres psql -U expense_user -d golang_expense_tracker \
  -c 'SELECT id, item, amount FROM public.expenses;'
```

On a fresh database, the query returns `(0 rows)`. If you see `relation "expenses" does not exist`, rerun the migration against `golang_expense_tracker`. If the database itself is missing and `expense_user` has permission to create databases, run this command first.

```bash
docker compose exec postgres psql -U expense_user -d postgres \
  -c 'CREATE DATABASE golang_expense_tracker;'
```

To inspect a startup failure, run `docker compose logs postgres`.

## 4. Run the Go server

```bash
go run .
```

The server runs at `http://127.0.0.1:8000`. You can run the current in-memory API independently of PostgreSQL. The database setup prepares the table for the PostgreSQL store integration. Until that integration is wired in, API requests use memory and expenses disappear when the Go server restarts.

## 5. Test the API

Open another terminal and run these commands in Bash or Git Bash.

**Check the server**

```bash
curl http://127.0.0.1:8000/health
```

**Create an expense**

```bash
curl -X POST http://127.0.0.1:8000/expenses -H 'Content-Type: application/json' -d '{"item":"Dinner","amount":100}'
```

**List expenses**

```bash
curl http://127.0.0.1:8000/expenses
```

**Get or update the expense**

Use the ID returned by the POST request. These examples use ID `1`.

```bash
curl http://127.0.0.1:8000/expenses/1

curl -X PATCH http://127.0.0.1:8000/expenses/1 -H 'Content-Type: application/json' -d '{"item":"Lunch","amount":120}'
```

Items cannot be blank, and amounts must be positive integers.

## 6. Stop the server and PostgreSQL

Press `Ctrl+C` in the Go server terminal. From the project root, stop PostgreSQL.

```bash
docker compose stop postgres
```

This keeps the PostgreSQL volume for the next run. Start PostgreSQL again with `docker compose up -d postgres`.
