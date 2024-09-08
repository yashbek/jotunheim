CREATE TABLE profiles (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    elo BIGSERIAL PRIMARY KEY,
    date_joined timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);