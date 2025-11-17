migrate -source file://infra/postgres/migrations -database postgres://localhost:5432/go_todo\?sslmode=disable up


migrate -path infra/postgres/migrations -database postgres://localhost:5432/go_todo\?sslmode=disable force VERSION


if error: Dirty database version 4. Fix and force version.

pg_dump -Fc -f go_todo_backup.dump go_todo

fix sql
migrate -source file://infra/postgres/migrations -database "postgres://localhost:5432/go_todo?sslmode=disable" force [number of last good]

migrate -source file://infra/postgres/migrations -database "postgres://localhost:5432/go_todo?sslmode=disable" up


Verify the migration succeeded
ini
psql -d go_todo -c "SELECT version, dirty FROM schema_migrations;"
psql -d go_todo -c "\d+ todolist"
psql -d go_todo -c "SELECT sequencename FROM pg_sequences WHERE sequencename = 'todolist_id_seq';"

Inspect partial objects (e.g., maybe the sequence was created but table not). If any partial objects exist you can drop them and re-run:
arduino
psql -d go_todo -c "DROP SEQUENCE IF EXISTS todolist_id_seq;"
psql -d go_todo -c "DROP TABLE IF EXISTS todolist;"
