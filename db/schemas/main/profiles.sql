CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    elo SERIAL NOT NULL,
    date_joined timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);