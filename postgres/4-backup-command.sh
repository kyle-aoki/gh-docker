PGPASSWORD="hunter2" pg_dump -h localhost -p 5432 -U postgres --column-inserts postgres > pg-backup.sql
