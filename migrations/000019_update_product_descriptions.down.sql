-- Revert description and cards_count for Chhatrapati Shivaji Maharaj
UPDATE products SET
    description = 'Master the legend of the Maratha Empire. This collector''s edition includes 52 beautifully illustrated cards covering the Coronation, Guerrilla Warfare tactics, Administration (Ashta Pradhan), and the architectural marvels of Sea Forts.',
    cards_count = 52
WHERE slug = 'shivaji-maratha';

-- Revert description and cards_count for A to Z Indian Heritage
UPDATE products SET
    description = 'A to Z Indian Heritage Cards — learn culture, monuments & traditions through illustrations.',
    cards_count = 52
WHERE slug = 'a-to-z-indian-heritage';

-- Revert description and cards_count for NEET Topper's Blueprint
UPDATE products SET
    description = 'NEET Topper''s Blueprint — high-yield Biology, Physics & Chemistry cards by toppers.',
    cards_count = 150
WHERE slug = 'neet-toppers-blueprint';
