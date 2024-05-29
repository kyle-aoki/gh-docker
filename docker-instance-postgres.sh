# back up

sudo apt install postgresql-client-common -y
sudo apt install postgresql-client -y

current_date=$(date '+%Y_%m_%d_%H_%M')

PGPASSWORD="hunter2" pg_dump -h localhost -p 5432 -U postgres --column-inserts postgres > pg-backup.sql

# operations

docker pull postgres

docker run \
-d \
-p 5432:5432 \
-e POSTGRES_PASSWORD=hunter2 \
-v /home/ubuntu/postgres-data:/var/lib/postgresql/data \
postgres
