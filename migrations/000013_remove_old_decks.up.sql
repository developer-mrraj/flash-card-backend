-- Remove Ramayana and UPSC decks
DELETE FROM products WHERE slug IN ('ramayana-epic', 'upsc-polity');
