-- A store cannot fulfil an order without knowing where to send it.
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS ship_name TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS ship_phone TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS ship_address TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS ship_city TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS ship_state TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS notes TEXT NOT NULL DEFAULT '';

ALTER TABLE orders
    ADD CONSTRAINT orders_status_check
    CHECK (status IN ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled'));

CREATE INDEX IF NOT EXISTS idx_orders_user_created ON orders(user_id, created_at DESC);
