CREATE TABLE IF NOT EXISTS quiz (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    current_answers TEXT NOT NULL,
    count_questions TEXT NOT NULL,
    type_questions TEXT NULL,
    difficulty_questions TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);