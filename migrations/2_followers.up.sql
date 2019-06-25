CREATE TABLE IF NOT EXISTS followers (
    follower text REFERENCES users (email),
    follows  text REFERENCES users (email)
);
