CREATE TABLE gifts
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    lookup_id TEXT    NOT NULL,
    gift_type TEXT    NOT NULL,
    redeemed  INTEGER NOT NULL,
    img_url    TEXT    NOT NULL
);

