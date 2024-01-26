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