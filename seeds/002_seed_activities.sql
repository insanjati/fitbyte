INSERT INTO activities (
    id,
    user_id,
    activity_type,
    done_at,
    duration_in_minutes,
    calories_burned
) VALUES

-- John's activities
(
    'a1b2c3d4-e5f6-7890-1234-567890abcdef',
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'Running',
    '2024-01-15 08:30:00+00',
    30,
    300
),
(
    'a1b2c3d4-e5f6-7890-1234-567890abcde0',
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'Cycling',
    '2024-01-14 18:00:00+00',
    45,
    400
),
(
    'a1b2c3d4-e5f6-7890-1234-567890abcde1',
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'Swimming',
    '2024-01-13 07:00:00+00',
    60,
    500
),

-- Jane's activities
(
    'b2c3d4e5-f6g7-8901-2345-678901bcdefg',
    'f47ac10b-58cc-4372-a567-0e02b2c3d480',
    'Yoga',
    '2024-01-15 06:00:00+00',
    60,
    200
),
(
    'b2c3d4e5-f6g7-8901-2345-678901bcdef0',
    'f47ac10b-58cc-4372-a567-0e02b2c3d480',
    'HIIT',
    '2024-01-14 17:30:00+00',
    20,
    250
),
(
    'b2c3d4e5-f6g7-8901-2345-678901bcdef1',
    'f47ac10b-58cc-4372-a567-0e02b2c3d480',
    'Dancing',
    '2024-01-12 19:00:00+00',
    40,
    180
),

-- Mike's activities
(
    'c3d4e5f6-g7h8-9012-3456-789012cdefgh',
    'f47ac10b-58cc-4372-a567-0e02b2c3d481',
    'Walking',
    '2024-01-15 07:30:00+00',
    25,
    120
),
(
    'c3d4e5f6-g7h8-9012-3456-789012cdefg0',
    'f47ac10b-58cc-4372-a567-0e02b2c3d481',
    'Stretching',
    '2024-01-14 20:00:00+00',
    15,
    50
),
(
    'c3d4e5f6-g7h8-9012-3456-789012cdefg1',
    'f47ac10b-58cc-4372-a567-0e02b2c3d481',
    'Hiking',
    '2024-01-13 09:00:00+00',
    90,
    350
);