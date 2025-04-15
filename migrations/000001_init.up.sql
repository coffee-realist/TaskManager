CREATE TYPE STATUS AS ENUM (
    'CREATED',
    'BOOKED',
    'COMPLETED'
);

CREATE TYPE ROLE AS ENUM (
    'WORKER',
    'PUBLISHER'
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    team VARCHAR(255) NOT NULL UNIQUE,
    role ROLE NOT NULL,
    password_hash VARCHAR(255) NOT NULl
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    project VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status STATUS NOT NULL,
    publisher INT REFERENCES,
    booked_by INT NOT NULL,
    booked_at TIMESTAMP DEFAULT NOW(),
    status_updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULl,
    expires_at TIMESTAMP NOT NULl,
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO users (username, team, role, password_hash)
VALUES ('w1', 't1', 'WORKER', '$2a$14$RXMhgf2Dw57a/QIfQdm9Je/eyHuP9bYnaN6Zh4VpnJOzPUI89nv3i')
    ON CONFLICT (username) DO NOTHING;
INSERT INTO users (username, team, role, password_hash)
VALUES ('w2', 't1', 'WORKER', '$2a$14$/qn/3eWeXeR0XE6IsTuGQunMvTLP6Q6iuYsWFWtuAoPH0IK4mXE3G')
    ON CONFLICT (username) DO NOTHING;
INSERT INTO users (username, team, role, password_hash)
VALUES ('w3', 't2', 'WORKER', '$2a$14$rAzHF.5905L8QuZZWFJUwuo5YPk/sMH8kPltEHMwIzEKDdNhfOHOK')
    ON CONFLICT (username) DO NOTHING;
INSERT INTO users (username, team, role, password_hash)
VALUES ('w4', 't2', 'WORKER', '$2a$14$2jIavKiNHAbaRe4Rw3Mh0OMNP0AuLDYNJU7xXU/Mws1vg/cnygoH2')
    ON CONFLICT (username) DO NOTHING;
INSERT INTO users (username, team, role, password_hash)
VALUES ('w5', 't2', 'WORKER', '$2a$14$wzp1nwmJWaXf6m68ozRamu5sA9ZnpUL6WYzILOzyA3SXB4T6ucNxS')
    ON CONFLICT (username) DO NOTHING;
INSERT INTO users (username, team, role, password_hash)
VALUES ('p1', 't1', 'PUBLISHER', '$2a$14$FPdZwLRGERY//aSKbvzw7Or3dhO2wjlg6m9nI6UuaIbP1shP4/zYC')
    ON CONFLICT (username) DO NOTHING;
INSERT INTO users (username, team, role, password_hash)
VALUES ('p2', 't2', 'PUBLISHER', '$2a$14$DAM4NAj4Xy.H8uOQsxa4b.q3GJ8U21d4lHusJ.nEp.Kzdn4tNj6y.')
    ON CONFLICT (username) DO NOTHING;
INSERT INTO users (username, team, role, password_hash)
VALUES ('p3', 't2', 'PUBLISHER', '$2a$14$fjZjgMT9jR2OEL1IsnEqKu1Zkp3kD/UPUeQV2wQ5khbtGFE2Kolpa')
    ON CONFLICT (username) DO NOTHING;
