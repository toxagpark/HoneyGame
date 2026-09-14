CREATE SCHEMA honey;

CREATE TABLE honey.users (
    id SERIAL PRIMARY KEY,
    tg_chat_id BIGINT NOT NULL UNIQUE,
    user_name VARCHAR(255) NOT NULL
);

CREATE TABLE honey.user_honey (
    user_id INTEGER PRIMARY KEY REFERENCES honey.users(id) ON DELETE CASCADE,
    honey BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE honey.challenges (
    id SERIAL PRIMARY KEY,
    amount BIGINT NOT NULL,
    winner_user_id INTEGER NOT NULL REFERENCES honey.users(id),
    loser_user_id INTEGER NOT NULL REFERENCES honey.users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_challenges_winner ON honey.challenges(winner_user_id);
CREATE INDEX idx_challenges_loser ON honey.challenges(loser_user_id);
CREATE INDEX idx_challenges_completed_at ON honey.challenges(completed_at);
