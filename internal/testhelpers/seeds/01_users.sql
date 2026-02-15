-- Seed test users
INSERT INTO users (id, email, first_name, last_name, created_at, updated_at)
VALUES
(
    '550e8400-e29b-41d4-a716-446655440001',
    'admin@test.com',
    'Test',
    'Admin',
    NOW(),
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440002',
    'user1@test.com',
    'Test',
    'User1',
    NOW(),
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440003',
    'user2@test.com',
    'Test',
    'User2',
    NOW(),
    NOW()
);
