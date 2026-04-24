-- Add deck-specific fields to the products table
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS slug           VARCHAR(255) UNIQUE,
    ADD COLUMN IF NOT EXISTS title          VARCHAR(500),
    ADD COLUMN IF NOT EXISTS badge          VARCHAR(100),
    ADD COLUMN IF NOT EXISTS badge_class    VARCHAR(100),
    ADD COLUMN IF NOT EXISTS rating         NUMERIC(3, 1) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reviews        INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS original_price BIGINT,
    ADD COLUMN IF NOT EXISTS discount       VARCHAR(50),
    ADD COLUMN IF NOT EXISTS main_image     VARCHAR(500),
    ADD COLUMN IF NOT EXISTS gallery_images TEXT[] DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS cards_count    INT DEFAULT 0;
