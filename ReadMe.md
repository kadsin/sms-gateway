# SMS Gateway

---

## Setup

Create a `.env` file from `.env.example` and fill it.

Install dependencies:

```bash
make init
```

Run the server:

```bash
make server
```

Run the queue worker:

```bash
make queue-worker topic=express.sms.pending
```

Run the balance worker:

```bash
make balance-worker consumers=25
```

---

## Testing

Show help for testing options:

```bash
make test:help
```

Run all tests:

```bash
make test
```

Optional flags:

-   `path=./pkg/utils` – test a specific path
-   `scope=server` – test a specific scope in `./tests/`
-   `filter=UserBalance` – run tests matching a filter
-   `verbose=t` – run tests with detailed output
-   `race=t` – run tests in race detection mode
-   `count=3` – run tests multiple times

Example:

```bash
make test verbose=t filter=UserHandler race=t
```

---

## Docker

```bash
make docker:up     # start containers
make docker:down   # stop containers
```

---

## Migrations

### Database Migrations

Create a new migration:

```bash
make migrate:db:create name="add_table"
```

Apply migrations:

```bash
make migrate:db
```

Rollback last migration:

```bash
make migrate:db:rollback
```

Check migration version:

```bash
make migrate:db:version
```

Check migration status:

```bash
make migrate:db:status
```

### Analytics Migrations

Create a new migration:

```bash
make migrate:analytics:create name="add_table"
```

Apply migrations:

```bash
make migrate:analytics
```

Rollback last migration:

```bash
make migrate:analytics:rollback
```

Check migration version:

```bash
make migrate:analytics:version
```

Check migration status:

```bash
make migrate:analytics:status
```

---

## Swagger Docs

-   **URL:** [http://localhost:3000/docs](http://localhost:3000/docs)
-   **Username:** `admin`
-   **Password:** `123456`

You can change the username and password in `.env`.
