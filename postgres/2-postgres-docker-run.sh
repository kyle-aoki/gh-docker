docker run \
-d \
-p 5432:5432 \
-e POSTGRES_PASSWORD=hunter2 \
-v /home/ubuntu/postgres-data:/var/lib/postgresql/data \
postgres
