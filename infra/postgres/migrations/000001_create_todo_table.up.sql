CREATE SEQUENCE IF NOT EXISTS todos_id_seq;

CREATE TABLE IF NOT EXISTS todos (
    id INTEGER NOT NULL DEFAULT nextval('todos_id_seq'),
    title VARCHAR(255) NOT NULL,
    done BOOL NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id)
);

ALTER SEQUENCE todos_id_seq OWNED BY todos.id;
