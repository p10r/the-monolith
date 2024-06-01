CREATE TABLE IF NOT EXISTS serve_matches
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    home_name  TEXT                              NOT NULL,
    away_name  TEXT                              NOT NULL,
    start_time INTEGER                           NOT NULL,
    country    TEXT                              NOT NULL,
    league     TEXT                              NOT NULL
);
