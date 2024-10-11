CREATE TABLE followers
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    user_id   INTEGER REFERENCES users (id),
    artist_id INTEGER REFERENCES artists (id)
);

