SELECT * FROM todos
WHERE
    user_id = :user_id
    AND
    list_id = :list_id
ORDER BY priority
