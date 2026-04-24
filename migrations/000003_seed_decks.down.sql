-- Remove seeded deck products by slug
DELETE FROM products WHERE slug IN ('shivaji-maratha', 'ramayana-epic', 'upsc-polity');
