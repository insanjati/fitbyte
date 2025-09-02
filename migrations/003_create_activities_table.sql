-- alternatif 
-- DROP TABLE IF EXISTS activities;
-- CREATE TYPE activity_enum AS ENUM ('Walking', 'Yoga', 'Stretching', 'Cycling', 'Swimming','Dancing', 'Hiking', 'Running', 'HIIT', 'JumpRope');

-- CREATE TABLE activities (
--     id UUID PRIMARY KEY,
--     user_id INT NOT NULL REFERENCES users(id) on DELETE CASCADE,
--     activity_type activity_enum NOT NULL,
--     done_at TIMESTAMPTZ NOT NULL,
--     duration_minutes INT NOT NULL CHECK (duration_minutes >= 1),
--     calories_burned INT NOT NULL,
--     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

DROP TABLE IF EXISTS activities;

CREATE TABLE activities (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    activity_type VARCHAR(50) NOT NULL CHECK (activity_type IN ('Walking', 'Yoga', 'Stretching', 'Cycling', 'Swimming', 'Dancing', 'Hiking', 'Running', 'HIIT', 'JumpRope')),
    done_at TIMESTAMP WITH TIME ZONE NOT NULL,
    duration_in_minutes INTEGER NOT NULL CHECK (duration_in_minutes > 0),
    calories_burned INTEGER NOT NULL CHECK (calories_burned > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- CREATE INDEX IF NOT EXISTS idx_activities_user_id ON activities(user_id);
-- CREATE INDEX IF NOT EXISTS idx_activities_done_at ON activities(done_at);
-- CREATE INDEX IF NOT EXISTS idx_activities_activity_type ON activities(activity_type);
-- CREATE INDEX IF NOT EXISTS idx_activities_calories_burned ON activities(calories_burned);
