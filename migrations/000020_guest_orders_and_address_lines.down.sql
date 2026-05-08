-- Rollback guest order and address line changes
ALTER TABLE user_addresses DROP COLUMN IF EXISTS address_line1;
ALTER TABLE user_addresses DROP COLUMN IF EXISTS address_line2;

ALTER TABLE orders DROP COLUMN IF EXISTS guest_name;
ALTER TABLE orders DROP COLUMN IF EXISTS guest_email;
ALTER TABLE orders DROP COLUMN IF EXISTS guest_phone;

DROP INDEX IF EXISTS idx_orders_guest_email;
DROP INDEX IF EXISTS idx_orders_guest_phone;

-- Note: We cannot re-add NOT NULL to user_id if there is existing data with NULL.
-- The original NOT NULL constraint is not restored in this rollback for safety.
