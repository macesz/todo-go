-- Rename the table
ALTER TABLE todolist RENAME TO todolists;

-- Rename the sequence
ALTER SEQUENCE todolist_id_seq RENAME TO todolists_id_seq;

-- Update the column default to use the new sequence name
ALTER TABLE todolists ALTER COLUMN id SET DEFAULT nextval('todolists_id_seq');

-- Update sequence ownership to the new table name
ALTER SEQUENCE todolists_id_seq OWNED BY todolists.id;

-- Rename the primary key constraint
ALTER TABLE todolists RENAME CONSTRAINT todolist_pkey TO todolists_pkey;
