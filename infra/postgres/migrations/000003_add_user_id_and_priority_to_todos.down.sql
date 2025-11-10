-- Remove priority column
ALTER TABLE todos
DROP COLUMN priority;

-- Remove user_id column
ALTER TABLE todos
DROP COLUMN user_id;
