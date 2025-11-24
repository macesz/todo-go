SELECT * FROM todos
WHERE
    user_id = :user_id
    AND
    todolist_id = :todolist_id
ORDER BY priority
