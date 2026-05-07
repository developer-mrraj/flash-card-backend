-- Force update images
UPDATE products SET main_image = '/images/a_z_800_600.png' WHERE slug = 'a-to-z-indian-heritage';
UPDATE products SET main_image = '/images/maharaj_800_600.png' WHERE slug = 'shivaji-maratha';
UPDATE products SET main_image = '/images/neet_800_600.png' WHERE slug = 'neet-toppers-blueprint';

-- Create a helper function to safely migrate IDs
CREATE OR REPLACE FUNCTION sync_product_id(target_slug TEXT, desired_id UUID)
RETURNS VOID AS $$
DECLARE
    current_id UUID;
BEGIN
    -- Get the current ID for the slug
    SELECT id INTO current_id FROM products WHERE slug = target_slug;

    -- If the product exists and the ID is different from desired
    IF current_id IS NOT NULL AND current_id != desired_id THEN
        
        -- Insert a duplicate product with the desired ID
        INSERT INTO products (id, name, description, price, stock_quantity, slug, title, badge, badge_class, rating, reviews, original_price, discount, main_image, gallery_images, cards_count, created_at, updated_at)
        SELECT desired_id, name, description, price, stock_quantity, target_slug || '-temp', title, badge, badge_class, rating, reviews, original_price, discount, main_image, gallery_images, cards_count, created_at, updated_at
        FROM products WHERE id = current_id;

        -- Update foreign keys in related tables
        UPDATE reviews SET product_id = desired_id WHERE product_id = current_id;
        
        -- If you have order_items or other tables, update them too.
        -- Assuming order_items table exists, we wrap in an exception block in case it doesn't.
        BEGIN
            EXECUTE 'UPDATE order_items SET product_id = $1 WHERE product_id = $2' USING desired_id, current_id;
        EXCEPTION WHEN undefined_table THEN
            -- do nothing
        END;

        -- Delete old product
        DELETE FROM products WHERE id = current_id;

        -- Fix the slug on the new product back to the original
        UPDATE products SET slug = target_slug WHERE id = desired_id;

    END IF;
END;
$$ LANGUAGE plpgsql;

-- Sync the IDs to match the hardcoded frontend values
SELECT sync_product_id('a-to-z-indian-heritage', 'e299a8e0-9bee-4911-8d5c-28497c324c44');
SELECT sync_product_id('shivaji-maratha', 'a4d312a9-b63e-4327-8b3f-dd42bfe738bf');
SELECT sync_product_id('neet-toppers-blueprint', '5d73cd3d-882a-49da-a163-8a8b0f73e309');
SELECT sync_product_id('ramayana-epic', '358c2b76-8025-4b0d-9b1d-91b5c490ec31');
SELECT sync_product_id('upsc-polity', '7b11d87f-9486-4f4f-80d5-5eb7127e2d93');

-- Drop the helper function
DROP FUNCTION sync_product_id(TEXT, UUID);
