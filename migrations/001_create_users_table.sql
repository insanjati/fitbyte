CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    preference VARCHAR(255),
    weight_unit VARCHAR(50),
    height_unit VARCHAR(50),
    weight DECIMAL(5,2),
    height DECIMAL(5,2),
    image_uri VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);