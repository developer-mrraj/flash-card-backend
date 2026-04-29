DROP TABLE IF EXISTS payment_history;
ALTER TABLE orders DROP COLUMN razorpay_order_id;
ALTER TABLE orders DROP COLUMN razorpay_payment_id;
