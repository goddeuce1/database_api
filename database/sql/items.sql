CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users CASCADE;
CREATE UNLOGGED TABLE IF NOT EXISTS users (
    nickname    CITEXT PRIMARY KEY COLLATE "POSIX",
    email       CITEXT UNIQUE NOT NULL,
    fullname    TEXT NOT NULL,
    about       TEXT
);

CREATE INDEX users_nickname ON users(nickname);

DROP TABLE IF EXISTS forums CASCADE;
CREATE UNLOGGED TABLE IF NOT EXISTS forums (
    slug        CITEXT PRIMARY KEY,
    posts       INTEGER DEFAULT 0,
    threads     INTEGER DEFAULT 0,
    title       TEXT NOT NULL,
    "user"      CITEXT NOT NULL
);

CREATE INDEX forums_slug ON forums(slug);

DROP TABLE IF EXISTS fu_table CASCADE;
CREATE UNLOGGED TABLE IF NOT EXISTS fu_table (
    nickname    CITEXT NOT NULL,
    forum       CITEXT NOT NULL,
    CONSTRAINT fu_table_constraint UNIQUE(nickname, forum)
);

CREATE INDEX fu_table_forum_nickname ON fu_table(forum, nickname);

DROP TABLE IF EXISTS threads CASCADE;
CREATE UNLOGGED TABLE IF NOT EXISTS threads (
    id          SERIAL PRIMARY KEY,
    slug        CITEXT UNIQUE,
    created     TIMESTAMPTZ(3)  DEFAULT now(),
    message     TEXT NOT NULL,
    title       TEXT NOT NULL,
    votes       INTEGER DEFAULT 0,
    forum       CITEXT NOT NULL,
    author      CITEXT NOT NULL
);

CREATE INDEX threads_forum_created ON threads(forum, created);

DROP TABLE IF EXISTS posts CASCADE;
CREATE UNLOGGED TABLE IF NOT EXISTS posts (
    id          SERIAL PRIMARY KEY,
    created     TIMESTAMPTZ(3) DEFAULT now(),
    isedited    BOOLEAN DEFAULT FALSE,
    message     TEXT NOT NULL,
    path        BIGINT[],
    author      CITEXT NOT NULL,
    forum       CITEXT NOT NULL,
    thread      INTEGER NOT NULL,
    parent      INTEGER DEFAULT 0
);

CREATE INDEX posts_thread_path ON posts(thread, path);
CREATE INDEX posts_thread_parent_id ON posts(thread, parent, id);
CREATE INDEX posts_thread_id ON posts(thread, id);

DROP TABLE IF EXISTS votes CASCADE;
CREATE UNLOGGED TABLE IF NOT EXISTS votes (
    voice       SMALLINT NOT NULL,
    nickname    CITEXT NOT NULL,
    thread      INTEGER NOT NULL,
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

CREATE OR REPLACE FUNCTION insert_path()
    RETURNS TRIGGER AS
    $$
    BEGIN
    IF new.parent IS NULL
        THEN new.path := (ARRAY [new.id]);
        RETURN new;
    END IF;
    new.path := (SELECT array_append(p.path, new.id::bigint) FROM posts p where p.id = new.parent);
        RETURN NEW;
    END;
    $$
    LANGUAGE 'plpgsql';

CREATE TRIGGER insert_vote
    AFTER INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE insert_vote();

CREATE TRIGGER check_vote
    AFTER UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE check_vote();

CREATE TRIGGER path
    BEFORE INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE insert_path();


