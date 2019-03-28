DROP TABLE IF EXISTS errors CASCADE;
CREATE TABLE IF NOT EXISTS errors (
    "message" TEXT
);

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE IF NOT EXISTS users (
    "about"     TEXT,
    "email"     CITEXT UNIQUE NOT NULL,
    "fullname"	CITEXT NOT NULL,
    "nickname"  CITEXT UNIQUE PRIMARY KEY NOT NULL
);

DROP TABLE IF EXISTS forums CASCADE;
CREATE TABLE IF NOT EXISTS forums (
    "posts"     INT DEFAULT 0,
    "slug"      CITEXT UNIQUE NOT NULL,
    "threads"   INT DEFAULT 0,
    "title"     TEXT NOT NULL,
    "user"      CITEXT NOT NULL REFERENCES users("nickname")
);

DROP TABLE IF EXISTS threads CASCADE;
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

DROP TABLE IF EXISTS posts CASCADE;
CREATE TABLE IF NOT EXISTS posts (
    "author"    CITEXT NOT NULL REFERENCES users("nickname"),
    "created"   TIMESTAMPTZ(3) DEFAULT now(),
    "forum"     CITEXT NOT NULL REFERENCES forums("slug"),
    "id"        SERIAL UNIQUE PRIMARY KEY NOT NULL,
    "isedited"  BOOLEAN DEFAULT FALSE,
    "message"   TEXT NOT NULL,
    "parent"    INT DEFAULT 0,
    "thread"    INT DEFAULT 0 REFERENCES threads("id"),
    "path" 	  	BIGINT []
);

DROP TABLE IF EXISTS votes CASCADE;
CREATE TABLE IF NOT EXISTS votes (
    "nickname"	CITEXT NOT NULL,
    "thread"	INT NOT NULL REFERENCES threads("id"),
    "voice" 	INT NOT NULL
);


