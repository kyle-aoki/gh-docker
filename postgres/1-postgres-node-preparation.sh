sudo apt update -y
sudo apt install docker.io -y

sudo chmod 666 /var/run/docker.sock

docker pull postgres
