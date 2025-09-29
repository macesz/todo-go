INSERT INTO todos (title)
VALUES (:title)
RETURNING id, created_at;
