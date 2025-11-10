-- Add user_id foreign key
ALTER TABLE todos
ADD COLUMN user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE;

-- Add priority column
ALTER TABLE todos
ADD COLUMN priority INTEGER NOT NULL DEFAULT 0;
