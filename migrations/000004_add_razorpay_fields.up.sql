ALTER TABLE orders ADD COLUMN razorpay_order_id VARCHAR(255);
ALTER TABLE orders ADD COLUMN razorpay_payment_id VARCHAR(255);

CREATE TABLE IF NOT EXISTS payment_history (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    razorpay_order_id VARCHAR(255) NOT NULL,
    razorpay_payment_id VARCHAR(255) UNIQUE NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(50) DEFAULT 'success',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
