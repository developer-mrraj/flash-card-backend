-- Revert deck-specific fields from the products table
ALTER TABLE products
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS title,
    DROP COLUMN IF EXISTS badge,
    DROP COLUMN IF EXISTS badge_class,
    DROP COLUMN IF EXISTS rating,
    DROP COLUMN IF EXISTS reviews,
    DROP COLUMN IF EXISTS original_price,
    DROP COLUMN IF EXISTS discount,
    DROP COLUMN IF EXISTS main_image,
    DROP COLUMN IF EXISTS gallery_images,
    DROP COLUMN IF EXISTS cards_count;
