
-- +migrate Up
INSERT INTO users (
    id,
    full_name,
    phone_number,
    password,
    is_active,
    email
) VALUES
(
    '550e8400-e29b-41d4-a716-446655440000',
    'Pratama',
    '081234567890',
    '$2a$10$examplehashedpassword1',
    1,
    'ammar.pratama@example.com'
),
(
    '550e8400-e29b-41d4-a716-446655440001',
    'Budi Santoso',
    '081298765432',
    '$2a$10$examplehashedpassword2',
    1,
    'budi.santoso@example.com'
),
(
    '550e8400-e29b-41d4-a716-446655440002',
    'Siti Rahma',
    '081277788899',
    '$2a$10$examplehashedpassword3',
    0,
    'siti.rahma@example.com'
);


INSERT INTO wallets (
    id,
    user_id,
    balance,
    currency,
    is_active
) VALUES
(
    '660e8400-e29b-41d4-a716-446655440000',
    '550e8400-e29b-41d4-a716-446655440000',
    1500000.00,
    'IDR',
    1
),
(
    '660e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440001',
    250000.00,
    'IDR',
    1
),
(
    '660e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440002',
    10000000.00,
    'IDR',
    0
);
-- +migrate Down
