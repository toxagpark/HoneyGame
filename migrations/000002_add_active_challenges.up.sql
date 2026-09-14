CREATE TABLE honey.active_challenges (
    id SERIAL PRIMARY KEY,
    creator_user_id INTEGER NOT NULL UNIQUE REFERENCES honey.users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_active_challenges_creator ON honey.active_challenges(creator_user_id);
