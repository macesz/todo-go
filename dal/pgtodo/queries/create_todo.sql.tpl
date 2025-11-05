INSERT INTO todos (user_id, title, priority)
VALUES (:user_id, :title, :priority)
RETURNING id, created_at;
