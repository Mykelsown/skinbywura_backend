CREATE UNIQUE INDEX IF NOT EXISTS products_name_key ON products (lower(name));
