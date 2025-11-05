UPDATE todos
SET title = :title, done = :done, priority = :priority
WHERE
    id = :id;
