UPDATE todos
SET title = :title, done = :done
WHERE
    id = :id;
