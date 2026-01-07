CREATE SEQUENCE IF NOT EXISTS todolist_id_seq;

CREATE TABLE IF NOT EXISTS todolist (
    id INTEGER NOT NULL DEFAULT nextval('todolist_id_seq'),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    color VARCHAR(255),
    labels VARCHAR(255),
    deleted BOOL NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id)
);

ALTER SEQUENCE todolist_id_seq OWNED BY todolist.id;
