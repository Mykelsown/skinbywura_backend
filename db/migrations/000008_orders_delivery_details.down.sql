ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_status_check;
DROP INDEX IF EXISTS idx_orders_user_created;
ALTER TABLE orders
    DROP COLUMN IF EXISTS ship_name,
    DROP COLUMN IF EXISTS ship_phone,
    DROP COLUMN IF EXISTS ship_address,
    DROP COLUMN IF EXISTS ship_city,
    DROP COLUMN IF EXISTS ship_state,
    DROP COLUMN IF EXISTS notes;
