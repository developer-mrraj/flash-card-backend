CREATE TABLE IF NOT EXISTS featured_collections (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    image VARCHAR(500) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO featured_collections (id, title, image, created_at, updated_at) VALUES
(gen_random_uuid(), 'EPIC SAGAS', 'http://localhost:8080/images/feat_history.png', NOW(), NOW()),
(gen_random_uuid(), 'EXAM PRO', 'http://localhost:8080/images/feat_upsc.png', NOW(), NOW()),
(gen_random_uuid(), 'TRADITIONS', 'http://localhost:8080/images/feat_mythology.png', NOW(), NOW());
