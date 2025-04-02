CREATE TABLE IF NOT EXISTS series (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) UNIQUE NOT NULL,
    status TEXT CHECK (status IN ('Plan to Watch', 'Watching', 'Dropped', 'Completed')) NOT NULL,
    last_episode_watched INT DEFAULT 0,
    total_episodes INT,
    ranking INT
);
