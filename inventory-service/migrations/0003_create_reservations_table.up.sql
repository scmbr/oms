CREATE TABLE  reservations(
    reservation_id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    product_id UUID NOT NULL REFERENCES products(product_id),
    quantity INTEGER NOT NULL CHECK (quantity>0),
    status VARCHAR(20) NOT NULL CHECK (status IN ('PENDING','RESERVED','FAILED','CANCELLED','EXPIRED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expired_at TIMESTAMPTZ 
);

CREATE INDEX idx_reservations_order_id ON reservations(order_id);
