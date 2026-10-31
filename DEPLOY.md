# Going live: step by step

Order matters. Do them top to bottom. Total time: about 45 minutes, plus waiting for DNS if you use a domain.

## 0. What you will end up with

- Storefront on Vercel (the React site).
- API and database on Render (this repo).
- Vercel forwards every `/api/...` request to Render. The browser only ever talks to the Vercel address, so login cookies work on every browser, including Safari and iPhones.

## 1. Email (Resend), 10 minutes

1. Create an account at resend.com.
2. Add and verify the shop's domain (Domains > Add). Without a verified domain Resend only delivers to your own address.
3. Create an API key (API Keys > Create). Keep it for step 2.
4. Decide the sender, for example `Skin By Wura <orders@yourdomain.com>`.

## 2. API and database on Render

1. Push this repo to GitHub (already done).
2. Render dashboard > New > Blueprint > pick this repo. Render reads `render.yaml` and creates the web service and Postgres database.
3. When asked, fill in the values marked "sync: false":

| Key | Value |
| --- | --- |
| `ALLOWED_ORIGINS` | the storefront address, e.g. `https://skinbywura.vercel.app` (no trailing slash). Add the real domain later. |
| `FRONTEND_URL` | same address |
| `RESEND_API_KEY` | key from step 1 |
| `MAIL_FROM` | sender from step 1 |
| `STORE_EMAIL` | inbox that gets contact messages and new order alerts |
| `PAYMENT_NOTE` | how customers pay, e.g. `Transfer to GTBank 0123456789, Skin By Wura. Use your order number as the reference.` |

4. Wait for the first deploy. Open `https://<your-service>.onrender.com/api/health`. You should see `{"status":"ok"}`.
5. Tables are created automatically at startup. Nothing else to run.

Note: the `starter` plan keeps the API awake. The free plan sleeps and the first request after a pause takes about a minute.

## 3. Add products

Option A, demo products: in the Render Postgres page copy the External Database URL, then
`psql "<that url>" -f db/seed.sql`.

Option B, real products (recommended): copy `catalog.example.csv` to `catalog.csv`, fill it in, then from a computer with Go installed:
`DATABASE_URL="<external url>" go run ./cmd/importproducts -file catalog.csv`.

Rules: prices are whole naira, `image_url` is a path such as `/images/serum.jpg` (put the image file in the storefront repo under `public/images/`), leave `rating` and `review_count` at 0 until you have real reviews.

## 4. Storefront on Vercel

1. Open `vercel.json` in the frontend repo and replace `REPLACE-WITH-YOUR-RENDER-SERVICE.onrender.com` with your Render address. Commit and push.
2. Vercel > Add New > Project > import the frontend repo. Framework: Vite (auto-detected).
3. Settings > Environment Variables, add the `VITE_` values from `.env.example` (contact email, WhatsApp, social links, payment note).
4. Deploy. Visit the Vercel address.
5. If `ALLOWED_ORIGINS` on Render was not the final address, update it and let Render redeploy.

## 5. Final checks (do all of them)

1. Sign up with a real email. You should land on your profile.
2. Add a product, open the bag, check out with delivery details. You get a confirmation email and `STORE_EMAIL` gets an alert.
3. Profile shows the order as "Awaiting confirmation".
4. Use "Forgot password". The email arrives, the link opens the reset page, the new password works.
5. Send a message from the Contact page. It arrives at `STORE_EMAIL`.
6. Refresh a deep link such as `/products/1`. It must load, not 404.

## 6. Day-to-day running

- New orders arrive by email. Mark progress in the database (Render > Postgres > Connect):
  `UPDATE orders SET status = 'shipped' WHERE id = 42;`
  Allowed: pending, confirmed, shipped, delivered, cancelled.
- See all orders: `SELECT id, status, total, ship_name, ship_phone, created_at FROM orders ORDER BY id DESC;`
- Contact messages are also saved: `SELECT * FROM contact_messages ORDER BY id DESC;`
- Render keeps database backups on paid plans. Confirm the plan you chose includes them.

## 7. Custom domain (optional)

Add the domain in Vercel, set its DNS as instructed, then update `ALLOWED_ORIGINS` and `FRONTEND_URL` on Render to the new address. The API address does not change.

## 8. Known limits (tell the client)

- No online card payment. Customers pay by transfer or on delivery as described in `PAYMENT_NOTE`, and the owner confirms by hand. The natural next build is Paystack.
- No admin dashboard. Orders and messages are handled by email plus the SQL above, and products by CSV import.
- Stock levels are not tracked.
- Rate limiting is per server instance, which is fine for one Render instance.
