CREATE TABLE articles (
    id serial primary key,
    title text,
    slug text,
    description text,
    body text,
    created timestamp,
    updated timestamp
);

