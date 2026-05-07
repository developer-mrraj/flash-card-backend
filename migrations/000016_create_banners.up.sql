CREATE TABLE IF NOT EXISTS banners (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot        VARCHAR(100) NOT NULL,
    title       VARCHAR(255) NOT NULL,
    subtitle    VARCHAR(500),
    cta_text    VARCHAR(100),
    cta_link    VARCHAR(500),
    image_url   VARCHAR(1000),
    bg_color    VARCHAR(255) DEFAULT '#0f0f23',
    text_color  VARCHAR(50)  DEFAULT '#ffffff',
    badge_text  VARCHAR(100),
    is_active   BOOLEAN      DEFAULT TRUE,
    sort_order  INT          DEFAULT 0,
    created_at  TIMESTAMPTZ  DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_banners_slot ON banners(slot);
CREATE INDEX IF NOT EXISTS idx_banners_is_active ON banners(is_active);

-- Seed 3 default banners (one per slot)
-- NOTE: Update image_url via PUT /api/admin/banners/{id} after uploading banners to Supabase storage
INSERT INTO banners (slot, title, subtitle, cta_text, cta_link, image_url, bg_color, text_color, badge_text, is_active, sort_order)
VALUES
(
    'home_mid',
    '🔥 Flash Sale — 25% Off All Decks!',
    'Master any topic faster with our hand-illustrated flashcards. Limited-time offer — grab yours before stock runs out!',
    'Shop Now',
    '/categories',
    '/25_off.png',
    'linear-gradient(135deg, #1a1a2e 0%, #16213e 60%, #0f3460 100%)',
    '#ffffff',
    'LIMITED TIME',
    TRUE,
    1
),
(
    'categories_top',
    '🎯 Exam Season is Here — Get the Edge!',
    'NEET • UPSC • Kids Learning — Top-rated decks trusted by 5,000+ students across India.',
    'Browse Decks',
    '/categories',
    '/banner_categories_top.png',
    'linear-gradient(135deg, #7c3aed 0%, #4f46e5 100%)',
    '#ffffff',
    'BESTSELLER',
    TRUE,
    1
);
