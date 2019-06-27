CREATE TABLE IF NOT EXISTS articles (
    id serial PRIMARY KEY,
    title text NOT NULL,
    slug text NOT NULL,
    description text,
    body text,
    created timestamptz NOT NULL DEFAULT NOW(),
    updated timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    email text PRIMARY KEY,
    name text UNIQUE NOT NULL,
    bio text,
    image text, -- base64
    password text NOT NULL
);
