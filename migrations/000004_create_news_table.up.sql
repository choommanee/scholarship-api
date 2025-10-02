CREATE TABLE IF NOT EXISTS news (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    summary TEXT NOT NULL,
    image_url VARCHAR(255),
    publish_date TIMESTAMP NOT NULL,
    expire_date TIMESTAMP,
    category VARCHAR(50) NOT NULL,
    tags TEXT[],
    is_published BOOLEAN NOT NULL DEFAULT false,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
