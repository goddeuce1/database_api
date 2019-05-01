CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE IF NOT EXISTS users (
    nickname    CITEXT PRIMARY KEY,
    email       CITEXT UNIQUE NOT NULL,
    fullname    TEXT NOT NULL,
    about       TEXT
);

DROP TABLE IF EXISTS forums CASCADE;
CREATE TABLE IF NOT EXISTS forums (
    slug        CITEXT PRIMARY KEY,
    posts       INTEGER DEFAULT 0,
    threads     INTEGER DEFAULT 0,
    title       TEXT NOT NULL,
    "user"      CITEXT NOT NULL REFERENCES users
);

DROP TABLE IF EXISTS threads CASCADE;
CREATE TABLE IF NOT EXISTS threads (
    id          SERIAL PRIMARY KEY,
    slug        CITEXT UNIQUE,
    created     TIMESTAMPTZ(3)  DEFAULT now(),
    message     TEXT NOT NULL,
    title       TEXT NOT NULL,
    votes       INTEGER DEFAULT 0,
    forum       CITEXT NOT NULL REFERENCES forums,
    author      CITEXT NOT NULL REFERENCES users
);

DROP TABLE IF EXISTS posts CASCADE;
CREATE TABLE IF NOT EXISTS posts (
    id          SERIAL PRIMARY KEY,
    created     TIMESTAMPTZ(3) DEFAULT now(),
    isedited    BOOLEAN DEFAULT FALSE,
    message     TEXT NOT NULL,
    path        BIGINT [],
    author      CITEXT NOT NULL REFERENCES users,
    forum       CITEXT NOT NULL REFERENCES forums,
    thread      INTEGER NOT NULL REFERENCES threads,
    parent      INTEGER DEFAULT 0
);

DROP TABLE IF EXISTS votes CASCADE;
CREATE TABLE IF NOT EXISTS votes (
    voice       SMALLINT NOT NULL,
    nickname    CITEXT REFERENCES users,
    thread      INTEGER REFERENCES threads,
    CONSTRAINT votes_constraint UNIQUE(thread, nickname)
);

CREATE OR REPLACE FUNCTION insert_vote()
    RETURNS TRIGGER AS 
    $$
    BEGIN
        UPDATE threads SET votes = votes + new.voice WHERE id = new.thread;
        RETURN new;
    END;
    $$
    LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION check_vote()
    RETURNS TRIGGER AS 
        $$
        BEGIN
            IF new.voice = -1 AND old.voice = 1
                THEN UPDATE threads SET votes = votes - 2 WHERE id = new.thread;
            END IF;
            IF new.voice = 1 AND old.voice = -1
                THEN UPDATE threads SET votes = votes + 2 WHERE id = new.thread;
            END IF;
            RETURN new;
        END;
        $$
        LANGUAGE 'plpgsql';

CREATE TRIGGER insert_vote
    AFTER INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE insert_vote();

CREATE TRIGGER check_vote
    AFTER UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE check_vote();

