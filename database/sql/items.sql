CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS users (
    nickname    CITEXT PRIMARY KEY,
    email       CITEXT UNIQUE NOT NULL,
    fullname    TEXT NOT NULL,
    about       TEXT
);

CREATE TABLE IF NOT EXISTS forums (
    slug        CITEXT PRIMARY KEY,
    posts       INTEGER DEFAULT 0,
    threads     INTEGER DEFAULT 0,
    title       TEXT NOT NULL,
    "user"      CITEXT NOT NULL REFERENCES users
);

CREATE TABLE IF NOT EXISTS threads (
    id          SERIAL PRIMARY KEY,
    slug        CITEXT,
    created     TIMESTAMPTZ(3)  DEFAULT now(),
    message     TEXT NOT NULL,
    title       TEXT NOT NULL,
    votes       INTEGER DEFAULT 0,
    forum       CITEXT NOT NULL REFERENCES forums,
    author      CITEXT NOT NULL REFERENCES users
);

CREATE TABLE IF NOT EXISTS posts (
    id          SERIAL PRIMARY KEY,
    created     TIMESTAMPTZ(3) DEFAULT now(),
    isedited    BOOLEAN DEFAULT FALSE,
    message     TEXT NOT NULL,
    parent      INTEGER DEFAULT 0,
    path        BIGINT [],
    author      CITEXT NOT NULL REFERENCES users,
    forum       CITEXT NOT NULL REFERENCES forums,
    thread      INTEGER DEFAULT 0 REFERENCES threads
);

CREATE TABLE IF NOT EXISTS votes (
    voice       SMALLINT NOT NULL,
    nickname    CITEXT REFERENCES users,
    thread      INTEGER REFERENCES threads
);
