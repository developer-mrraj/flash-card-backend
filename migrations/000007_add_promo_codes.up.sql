CREATE TABLE IF NOT EXISTS promo_codes (
    id UUID PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    discount_percentage INTEGER NOT NULL CHECK (discount_percentage > 0 AND discount_percentage <= 100),
    is_active BOOLEAN DEFAULT TRUE,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Seed default promo codes
INSERT INTO promo_codes (id, code, discount_percentage, is_active) VALUES
  (gen_random_uuid(), 'FLASH20', 20, true),
  (gen_random_uuid(), 'WELCOME10', 10, true),
  (gen_random_uuid(), 'STUDY15', 15, true)
ON CONFLICT (code) DO NOTHING;
