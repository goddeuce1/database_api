CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS users (
    "about"     TEXT,
    "email"     CITEXT UNIQUE NOT NULL,
    "fullname"	TEXT NOT NULL,
    "nickname"  CITEXT UNIQUE PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS forums (
    "posts"     INTEGER DEFAULT 0,
    "slug"      CITEXT UNIQUE NOT NULL,
    "threads"   INTEGER DEFAULT 0,
    "title"     TEXT NOT NULL,
    "user"      CITEXT NOT NULL REFERENCES users("nickname")
);

CREATE TABLE IF NOT EXISTS threads (
    "author"    CITEXT NOT NULL REFERENCES users("nickname"),
    "created"   TIMESTAMPTZ(3)  DEFAULT now(),
    "forum"     CITEXT NOT NULL REFERENCES forums("slug"),
    "id"        SERIAL UNIQUE PRIMARY KEY NOT NULL,
    "message"   TEXT NOT NULL,
    "slug"      CITEXT,
    "title"     TEXT NOT NULL,
    "votes"     INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS posts (
    "author"    CITEXT NOT NULL REFERENCES users("nickname"),
    "created"   TIMESTAMPTZ(3) DEFAULT now(),
    "forum"     CITEXT NOT NULL REFERENCES forums("slug"),
    "id"        SERIAL UNIQUE PRIMARY KEY,
    "isedited"  BOOLEAN DEFAULT FALSE,
    "message"   TEXT NOT NULL,
    "parent"    INTEGER DEFAULT 0,
    "thread"    INTEGER DEFAULT 0 REFERENCES threads("id"),
    "path" 	  	BIGINT []
);

CREATE TABLE IF NOT EXISTS votes (
    "nickname"	CITEXT NOT NULL,
    "thread"	INTEGER NOT NULL REFERENCES threads("id"),
    "voice" 	INTEGER NOT NULL
);
