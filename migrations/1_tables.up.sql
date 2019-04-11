CREATE TABLE IF NOT EXISTS articles (
    id serial PRIMARY KEY,
    title text NOT NULL,
    slug text NOT NULL,
    description text,
    body text,
    created timestamptz NOT NULL DEFAULT NOW(),
    updated timestamptz NOT NULL DEFAULT NOW()
);

