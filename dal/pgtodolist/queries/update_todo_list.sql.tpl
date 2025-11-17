UPDATE todolist
SET title = :title, color = :color, labels = :labels
WHERE
    id = :id;
