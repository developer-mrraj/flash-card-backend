-- Update description and cards_count for Chhatrapati Shivaji Maharaj
UPDATE products SET
    description = 'Master the legend of the Maratha Empire. This collector''s edition includes 25 × 2 beautifully illustrated cards covering the Coronation, Guerrilla Warfare tactics, Administration (Ashta Pradhan), and the architectural marvels of Sea Forts.',
    cards_count = 50
WHERE slug = 'shivaji-maratha';

-- Update description and cards_count for A to Z Indian Heritage
UPDATE products SET
    description = 'Explore India''s rich cultural legacy from A to Z! 25 × 2 beautifully illustrated cards covering monuments, festivals, art forms, folk traditions, and heritage sites across every letter of the alphabet. Perfect for kids and curious learners.',
    cards_count = 50
WHERE slug = 'a-to-z-indian-heritage';

-- Update description and cards_count for NEET Topper's Blueprint
UPDATE products SET
    description = 'NEET Roadmap Cards — a complete NCERT-aligned study roadmap showing what to study, how to study, and when to revise for Biology, Physics & Chemistry like top NEET rankers.',
    cards_count = 50
WHERE slug = 'neet-toppers-blueprint';
