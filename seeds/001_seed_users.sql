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
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
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
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
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
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'yoga',
    'kg',
    'cm',
    80,
    175,
    null
);