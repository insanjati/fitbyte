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
    '20354d7a-e4fe-47af-8ff6-187bca92f3f9',
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'Running',
    '2024-01-15 08:30:00+00',
    30,
    300
),
(
    '20354d7a-e4fe-47af-8ff6-187bca92f3f8',
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'Cycling',
    '2024-01-14 18:00:00+00',
    45,
    400
),
(
    '20354d7a-e4fe-47af-8ff6-187bca92f3f7',
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'Swimming',
    '2024-01-13 07:00:00+00',
    60,
    500
),

-- Jane's activities
(
    '20354d7a-e4fe-47af-8ff6-187bca92f3f6',
    'f47ac10b-58cc-4372-a567-0e02b2c3d480',
    'Yoga',
    '2024-01-15 06:00:00+00',
    60,
    200
),
(
    '20354d7a-e4fe-47af-8ff6-187bca92f3f5',
    'f47ac10b-58cc-4372-a567-0e02b2c3d480',
    'HIIT',
    '2024-01-14 17:30:00+00',
    20,
    250
),
(
    '20354d7a-e4fe-47af-8ff6-187bca92f3f4',
    'f47ac10b-58cc-4372-a567-0e02b2c3d480',
    'Dancing',
    '2024-01-12 19:00:00+00',
    40,
    180
),

-- Mike's activities
(
    '20354d7a-e4fe-47af-8ff6-187bca92f3f3',
    'f47ac10b-58cc-4372-a567-0e02b2c3d481',
    'Walking',
    '2024-01-15 07:30:00+00',
    25,
    120
),
(
    '20354d7a-e4fe-47af-8ff6-187bca92f3f2',
    'f47ac10b-58cc-4372-a567-0e02b2c3d481',
    'Stretching',
    '2024-01-14 20:00:00+00',
    15,
    50
),
(
    '20354d7a-e4fe-47af-8ff6-187bca92f3f1',
    'f47ac10b-58cc-4372-a567-0e02b2c3d481',
    'Hiking',
    '2024-01-13 09:00:00+00',
    90,
    350
);