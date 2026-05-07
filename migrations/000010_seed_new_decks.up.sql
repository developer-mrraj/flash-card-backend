INSERT INTO products (
    id, name, description, price, stock_quantity, slug, title, badge, badge_class, rating, reviews, main_image, gallery_images, cards_count, created_at, updated_at
) VALUES
(
    'e299a8e0-9bee-4911-8d5c-28497c324c44',
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
    '5d73cd3d-882a-49da-a163-8a8b0f73e309',
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
