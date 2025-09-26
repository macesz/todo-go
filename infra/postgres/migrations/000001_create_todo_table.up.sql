CREATE TABLE IF NOT EXISTS todos(
   id serial PRIMARY KEY,
   title VARCHAR (255) NOT NULL,
   done BOOL NOT NULL DEFAULT false,
   created_at timestamp
);
