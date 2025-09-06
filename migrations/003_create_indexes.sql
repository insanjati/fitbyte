-- Users table indexes
CREATE INDEX idx_users_email ON users(email); -- For email lookup during authentication

-- Activities table indexes  
CREATE INDEX idx_activities_user_activity ON activities(user_id, id); -- For activity ownership checking
CREATE INDEX idx_activities_user_done_at ON activities(user_id, done_at DESC); -- For time-based filtering
CREATE INDEX idx_activities_user_type ON activities(user_id, activity_type); -- For activity type filtering
-- CREATE INDEX idx_activities_user_calories ON activities(user_id, calories_burned); -- For calories filtering (optional)