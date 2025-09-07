INSERT INTO users (
    id, 
    name, 
    email, 
    password, 
    preference, 
    weightUnit, 
    heightUnit, 
    weight, 
    height, 
    imageUri
) VALUES
(
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'John Doe',
    'john@example.com',
    '$2a$10$Jy7S4Ea8VgcbeVLrhNaQz.TdHKUEyMFzTkAercFa3TPqee/bKV67.', -- password123
    'cardio',
    'kg',
    'cm',
    75,
    180,
    'https://via.placeholder.com/150'
),
(
    'f47ac10b-58cc-4372-a567-0e02b2c3d480',
    'Jane Smith',
    'jane@example.com',
    '$2a$10$Jy7S4Ea8VgcbeVLrhNaQz.TdHKUEyMFzTkAercFa3TPqee/bKV67.',
    'strength',
    'lbs',
    'ft',
    140,
    66,
    'https://via.placeholder.com/150'
),
(
    'f47ac10b-58cc-4372-a567-0e02b2c3d481',
    'Mike Johnson',
    'mike@example.com',
    '$2a$10$Jy7S4Ea8VgcbeVLrhNaQz.TdHKUEyMFzTkAercFa3TPqee/bKV67.',
    'yoga',
    'kg',
    'cm',
    80,
    175,
    null
);