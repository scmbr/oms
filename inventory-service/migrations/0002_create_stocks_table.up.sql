CREATE TABLE stocks(
    product_id UUID PRIMARY KEY REFERENCES products(product_id),
    available INTEGER NOT NULL CHECK (available>=0)
);
