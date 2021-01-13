CREATE TABLE IF NOT EXISTS followers (
    follower int REFERENCES users (id),
    follows  int REFERENCES users (id),
    UNIQUE(follower, follows)
);
