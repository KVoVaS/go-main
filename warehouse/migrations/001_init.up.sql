CREATE TABLE products (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    stock INT NOT NULL DEFAULT 0
);

CREATE TABLE orders (
    id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL REFERENCES products(id),
    quantity INT NOT NULL,
    status TEXT NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- предзаполним один товар для теста
INSERT INTO products (id, name, stock) VALUES ('prod-1', 'Widget', 100);