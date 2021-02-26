CREATE TABLE IF NOT EXISTS followers (
    follower int REFERENCES users (id),
    followee int REFERENCES users (id),
    UNIQUE(follower, followee)
);
