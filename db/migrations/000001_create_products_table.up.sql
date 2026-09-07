CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price INTEGER NOT NULL,
    compare_at_price INTEGER,
    rating NUMERIC(3,1) NOT NULL CHECK (rating >= 0 AND rating <= 5),
    review_count INTEGER NOT NULL DEFAULT 0,
    skin_type TEXT NOT NULL,
    badge TEXT,
    image_url TEXT NOT NULL,
    description TEXT NOT NULL,
    volume TEXT NOT NULL
);
