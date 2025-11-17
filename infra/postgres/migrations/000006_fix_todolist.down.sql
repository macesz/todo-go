-- Rename primary key constraint back
ALTER TABLE todolists RENAME CONSTRAINT todolists_pkey TO todolist_pkey;

-- Rename sequence back
ALTER SEQUENCE todolists_id_seq RENAME TO todolist_id_seq;

-- Update column default back to old sequence name
ALTER TABLE todolists ALTER COLUMN id SET DEFAULT nextval('todolist_id_seq');

-- Update sequence ownership back
ALTER SEQUENCE todolist_id_seq OWNED BY todolists.id;

-- Rename table back
ALTER TABLE todolists RENAME TO todolist;
