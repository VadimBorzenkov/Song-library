CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(40) NOT NULL,
    song_name VARCHAR(50) NOT NULL,
    release_date DATE,
    text TEXT,
    link TEXT
);