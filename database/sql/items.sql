CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS users (
    nickname    CITEXT UNIQUE PRIMARY KEY NOT NULL,
    email       CITEXT UNIQUE NOT NULL,
    fullname    TEXT NOT NULL,
    about       TEXT
);

CREATE TABLE IF NOT EXISTS forums (
    slug        CITEXT UNIQUE NOT NULL,
    posts       INTEGER DEFAULT 0,
    threads     INTEGER DEFAULT 0,
    title       TEXT NOT NULL,
    "user"      CITEXT NOT NULL REFERENCES users("nickname")
);

CREATE TABLE IF NOT EXISTS threads (
    id          SERIAL UNIQUE PRIMARY KEY NOT NULL,
    slug        CITEXT,
    created     TIMESTAMPTZ(3)  DEFAULT now(),
    message     TEXT NOT NULL,
    title       TEXT NOT NULL,
    votes       INTEGER DEFAULT 0,
    forum       CITEXT NOT NULL REFERENCES forums("slug"),
    author      CITEXT NOT NULL REFERENCES users("nickname")
);

CREATE TABLE IF NOT EXISTS posts (
    id          SERIAL UNIQUE PRIMARY KEY,
    created     TIMESTAMPTZ(3) DEFAULT now(),
    isedited    BOOLEAN DEFAULT FALSE,
    message     TEXT NOT NULL,
    parent      INTEGER DEFAULT 0,
    path 	    BIGINT [],
    author      CITEXT NOT NULL REFERENCES users("nickname"),
    forum       CITEXT NOT NULL REFERENCES forums("slug"),
    thread      INTEGER DEFAULT 0 REFERENCES threads("id")
);

CREATE TABLE IF NOT EXISTS votes (
    voice 	    INTEGER NOT NULL,
    nickname    CITEXT NOT NULL REFERENCES users("nickname"),
    thread	    INTEGER NOT NULL REFERENCES threads("id")
);
