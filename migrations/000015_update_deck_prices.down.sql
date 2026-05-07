-- Revert prices to previous state
UPDATE products SET price = 999 WHERE slug = 'a-to-z-indian-heritage';
UPDATE products SET price = 1499 WHERE slug = 'neet-toppers-blueprint';
UPDATE products SET price = 1299 WHERE slug = 'shivaji-maratha';
