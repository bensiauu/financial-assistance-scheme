
INSERT INTO administrators (id, name, email, password_hash, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Root User',
    'root@root.com',
    crypt('root', gen_salt('bf')),
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;
