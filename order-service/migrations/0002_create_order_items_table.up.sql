CREATE TABLE order_items (
    item_id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    quantity INT NOT NULL,
    price NUMERIC(10,2) NOT NULL
);
