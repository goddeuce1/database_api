DROP TABLE IF EXISTS errors CASCADE;
CREATE TABLE IF NOT EXISTS errors (
    message TEXT
);

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE IF NOT EXISTS users (
    about     TEXT,
    email     VARCHAR(50) UNIQUE NOT NULL,
    fullname  TEXT NOT NULL,
    nickname  VARCHAR(50) UNIQUE PRIMARY KEY NOT NULL
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

DROP TABLE IF EXISTS forums CASCADE;
CREATE TABLE IF NOT EXISTS forums (
    posts     INT DEFAULT 0,
    slug      VARCHAR(50) UNIQUE NOT NULL,
    threads   INT DEFAULT 0,
    title     TEXT NOT NULL,
    user      TEXT NOT NULL REFERENCES users(nickname)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

DROP TABLE IF EXISTS threads CASCADE;
CREATE TABLE IF NOT EXISTS threads (
    author    TEXT NOT NULL REFERENCES users(nickname),
    created   TIMESTAMP(3)  DEFAULT CURRENT_TIMESTAMP(3),
    forum     TEXT NOT NULL REFERENCES forums(slug),
    id        SERIAL UNIQUE PRIMARY KEY NOT NULL,
    message   TEXT NOT NULL,
    slug      TEXT,
    title     TEXT NOT NULL,
    votes     INT DEFAULT 0
) ENGINE = InnoDB AUTO_INCREMENT = 42 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

DROP TABLE IF EXISTS posts CASCADE;
CREATE TABLE IF NOT EXISTS posts (
    author    TEXT NOT NULL REFERENCES users(nickname),
    created   TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    forum     TEXT NOT NULL REFERENCES forums(slug),
    id        SERIAL UNIQUE PRIMARY KEY NOT NULL,
    isEdited  BOOLEAN DEFAULT FALSE,
    message   TEXT NOT NULL,
    parent    INT,
    thread    INT DEFAULT 0 REFERENCES threads(id)
) ENGINE = InnoDB AUTO_INCREMENT = 42 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

DROP TABLE IF EXISTS votes CASCADE;
CREATE TABLE IF NOT EXISTS votes (
    nickname TEXT NOT NULL,
    slug 	 TEXT NOT NULL,
    slugid	 INT NOT NULL,
    voice 	 INT NOT NULL
);

