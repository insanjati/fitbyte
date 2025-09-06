CREATE TABLE activities (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    activity_type VARCHAR(50) NOT NULL CHECK (activity_type IN ('Walking', 'Yoga', 'Stretching', 'Cycling', 'Swimming', 'Dancing', 'Hiking', 'Running', 'HIIT', 'JumpRope')),
    done_at TIMESTAMP WITH TIME ZONE NOT NULL,
    duration_in_minutes INTEGER NOT NULL CHECK (duration_in_minutes > 0),
    calories_burned INTEGER NOT NULL CHECK (calories_burned > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);