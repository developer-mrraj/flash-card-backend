-- 1. Make user_id nullable on orders (to support guest orders)
ALTER TABLE orders ALTER COLUMN user_id DROP NOT NULL;

-- 2. Add guest information columns to orders
ALTER TABLE orders ADD COLUMN IF NOT EXISTS guest_name  VARCHAR(255);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS guest_email VARCHAR(255);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS guest_phone VARCHAR(20);

-- 3. Add address_line1 and address_line2 to user_addresses
ALTER TABLE user_addresses ADD COLUMN IF NOT EXISTS address_line1 TEXT;
ALTER TABLE user_addresses ADD COLUMN IF NOT EXISTS address_line2 TEXT;

-- 4. Migrate existing 'address' values into address_line1
UPDATE user_addresses SET address_line1 = address WHERE address_line1 IS NULL AND address IS NOT NULL;

-- 5. Index on guest_email for fast lookup during signup linking
CREATE INDEX IF NOT EXISTS idx_orders_guest_email ON orders (guest_email) WHERE guest_email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_orders_guest_phone ON orders (guest_phone) WHERE guest_phone IS NOT NULL;
