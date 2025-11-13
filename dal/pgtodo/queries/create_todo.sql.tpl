INSERT INTO todos (user_id, list_id title, priority, created_at)
VALUES (:user_id, :lst_id, :title, :priority, :created_at)
RETURNING id;
