-- Remove features JSON column from products table
ALTER TABLE products
    DROP COLUMN IF EXISTS features;
