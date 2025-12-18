CREATE TABLE products (
    product_id UUID PRIMARY KEY,
    title VARCHAR(255),
    sku VARCHAR(20) NOT NULL UNIQUE,
    price NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
