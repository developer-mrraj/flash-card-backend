INSERT INTO products (
    id, name, description, price, stock_quantity, slug, title, badge, badge_class, rating, reviews, main_image, gallery_images, cards_count, created_at, updated_at
) VALUES
(
    gen_random_uuid(),
    'A to Z Indian Heritage Cards',
    'A to Z Indian Heritage Cards — learn culture, monuments & traditions through illustrations.',
    999,
    100,
    'a-to-z-indian-heritage',
    'A to Z Indian Heritage Cards',
    'KIDS',
    'badge-info',
    4.9,
    1200,
    '/images/a_z_800_600.png',
    ARRAY[]::TEXT[],
    52,
    NOW(),
    NOW()
),
(
    gen_random_uuid(),
    'NEET Topper''s Blueprint',
    'NEET Topper''s Blueprint — high-yield Biology, Physics & Chemistry cards by toppers.',
    1499,
    100,
    'neet-toppers-blueprint',
    'NEET Topper''s Blueprint',
    'MUST HAVE',
    'badge-warning',
    4.8,
    3400,
    '/images/neet_800_600.png',
    ARRAY[]::TEXT[],
    150,
    NOW(),
    NOW()
);
