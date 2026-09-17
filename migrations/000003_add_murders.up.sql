CREATE TABLE honey.murders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES honey.users(id) ON DELETE CASCADE,
    hives INTEGER NOT NULL CHECK (hives > 0),
    amount BIGINT NOT NULL CHECK (amount > 0),
    honey BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_murders_user_id ON honey.murders(user_id);
CREATE INDEX idx_murders_created_at ON honey.murders(created_at);
