CREATE TABLE IF NOT EXISTS artists
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    ra_id          INTEGER                           NOT NULL,
    ra_slug        TEXT                              NOT NULL,
    name           TEXT                              NOT NULL,
    followers      TEXT                              NOT NULL,
    tracked_events TEXT                              NOT NULL
);

CREATE TABLE IF NOT EXISTS monitoring_events
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    event_type TEXT                              NOT NULL,
    data       TEXT                              NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    is_admin       INTEGER                           NOT NULL,
    telegram_id    INTEGER                           NOT NULL,
    tracked_events TEXT
);

CREATE TABLE IF NOT EXISTS serve_matches
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    home_name  TEXT                              NOT NULL,
    away_name  TEXT                              NOT NULL,
    start_time INTEGER                           NOT NULL,
    country    TEXT                              NOT NULL,
    league     TEXT                              NOT NULL
);
