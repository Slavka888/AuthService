INSERT INTO
    apps (id, name, secret)
VALUES
    ('1', 'App One', 'secret1') ON conflict DO nothing;