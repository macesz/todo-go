INSERT INTO todos (user_id, title, priority, created_at)
VALUES (:user_id, :title, :priority, :created_at)
RETURNING id;
