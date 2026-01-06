UPDATE todolists
SET title = :title, color = :color, labels = :labels, deleted = :deleted
WHERE
    id = :id;
