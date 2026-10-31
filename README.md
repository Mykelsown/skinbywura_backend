# Skin By Wura API

Go + PostgreSQL backend for the Skin By Wura storefront.

## What it does

| Area | Endpoints |
| --- | --- |
| Catalogue | `GET /api/products`, `GET /api/products/{id}` |
| Accounts | `POST /api/auth/signup`, `login`, `logout`, `forgot-password`, `reset-password`, `change-password`; `GET` and `PATCH /api/auth/me` |
| Cart and wishlist | `GET` and `PUT /api/cart`, `GET` and `PUT /api/wishlist` |
| Orders | `POST /api/orders` (needs delivery details), `GET /api/orders` |
| Contact | `POST /api/contact` |
| Health | `GET /api/health` |

Orders are created as `pending`. There is no online payment: the confirmation email carries the
`PAYMENT_NOTE` text (for example bank transfer details) and the owner confirms payment by hand.

## Run locally

1. Install Go 1.26+ and PostgreSQL. Create a database and user.
2. `cp .env.example .env` and set `DATABASE_URL`.
3. `make run`. Migrations apply automatically on start.
4. Optional demo catalogue: `make seed`.

With no `RESEND_API_KEY`, emails (including password reset links) are printed to the log.

## Configuration

See `.env.example`. In production (`APP_ENV=production`) the app refuses to start unless cookies are
`Secure` and `ALLOWED_ORIGINS` lists real origins.

## Managing products (no admin panel)

Copy `catalog.example.csv` to `catalog.csv`, edit it in a spreadsheet, then run
`make import` (locally) or `/app/importproducts -file catalog.csv` (inside the container).
Products are matched by name: existing rows update, new rows are created. Use `-dry-run` to validate.

Reviews and ratings: leave `rating` and `review_count` at 0 for products without real reviews.

## Order status

Update status in the database (`pending`, `confirmed`, `shipped`, `delivered`, `cancelled`):

```
UPDATE orders SET status = 'confirmed' WHERE id = 42;
```

## Deploy

See `DEPLOY.md`.

## Tests

`make test`. CI (`.github/workflows/ci.yml`) also boots the API against a real Postgres.
