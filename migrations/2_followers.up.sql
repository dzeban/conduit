CREATE TABLE IF NOT EXISTS followers (
    follower text REFERENCES users (name),
    follows  text REFERENCES users (name),
    UNIQUE(follower, follows)
);
