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

---

## Testing

```bash
make test
```

Optional flags:

* `path=./pkg/utils`
* `scope=api`
* `filter=Login`
* `verbose=true`
* `race=true`
* `count=3`

Example:

```bash
make test verbose=true filter=UserHandler
```

---

## Docker

```bash
make docker:up     # start containers
make docker:down   # stop containers
```

---

## Migrations

```bash
make migrate:create name="add_table"
make migrate        # apply
make migrate:rollback
make migrate:status
```

---

## Swagger Docs

* **URL:** [http://localhost:3000/docs](http://localhost:3000/docs)
* **Username:** `admin`
* **Password:** `123456`

Note that you can change username and password in `.env`.
