# Go Expense REST API

A simple API to create, list, and update expenses.

## 1. Install Go

Install Go 1.27.3 or newer from https://go.dev/dl/.

Check the installation:

```bash
go version
```

## 2. Download the project

Replace `<repository-url>` with this project's Git URL:

```bash
git clone <repository-url> go-lang-restapi
cd go-lang-restapi
```

If you downloaded the parent learning repository, open its `go-lang-restapi` folder instead.

## 3. Run the server

```bash
go run .
```

The server runs at `http://127.0.0.1:8000`.

## 4. Test the API

Open another terminal. Run these commands in Bash or Git Bash.

Check the server:

```bash
curl http://127.0.0.1:8000/health
```

Create an expense:

```bash
curl -X POST http://127.0.0.1:8000/expenses -H 'Content-Type: application/json' -d '{"item":"Dinner","amount":100}'
```

List expenses:

```bash
curl http://127.0.0.1:8000/expenses
```

Get expense 1:

```bash
curl http://127.0.0.1:8000/expenses/1
```

Update expense 1:

```bash
curl -X PATCH http://127.0.0.1:8000/expenses/1 -H 'Content-Type: application/json' -d '{"item":"Lunch","amount":120}'
```

## 5. Stop the server

Press `Ctrl+C` in the server terminal.

Expenses are stored in memory and disappear when the server restarts. Create an expense before trying to fetch or update it. Items cannot be blank, and amounts must be positive integers.
