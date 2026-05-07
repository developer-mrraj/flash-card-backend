-- Update price of all flashcard decks to 249 INR
-- Note: Price is stored in paise (cents), so 249 * 100 = 24900

UPDATE products 
SET price = 24900,
    original_price = 49900,
    discount = '50% OFF';
