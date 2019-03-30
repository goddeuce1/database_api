CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    "about"     TEXT,
    "email"     CITEXT UNIQUE NOT NULL,
    "fullname"	CITEXT NOT NULL,
    "nickname"  CITEXT UNIQUE PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS forums (
    "posts"     INT DEFAULT 0,
    "slug"      CITEXT UNIQUE NOT NULL,
    "threads"   INT DEFAULT 0,
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
    "votes"     INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS posts (
    "author"    CITEXT NOT NULL REFERENCES users("nickname"),
    "created"   TIMESTAMPTZ(3) DEFAULT now(),
    "forum"     CITEXT NOT NULL REFERENCES forums("slug"),
    "id"        SERIAL UNIQUE PRIMARY KEY NOT NULL,
    "isedited"  BOOLEAN DEFAULT FALSE,
    "message"   TEXT NOT NULL,
    "parent"    INT DEFAULT 0,
    "thread"    INT DEFAULT 0 REFERENCES threads("id"),
    "path" 	  	INT []
);

CREATE TABLE IF NOT EXISTS votes (
    "nickname"	CITEXT NOT NULL,
    "thread"	INT NOT NULL REFERENCES threads("id"),
    "voice" 	INT NOT NULL
);
