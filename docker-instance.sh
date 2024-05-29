###################################################################################################
#### install ######################################################################################
###################################################################################################
sudo apt update -y
sudo apt install docker.io -y
sudo chmod 666 /var/run/docker.sock
sudo apt install nginx -y

###################################################################################################
#### files ########################################################################################
###################################################################################################
SSH_HOST=dev
scp ./nginx.conf $SSH_HOST:nginx.conf
ssh $SSH_HOST sudo cp nginx.conf /etc/nginx/nginx.conf
ssh $SSH_HOST sudo nginx -s reload

GOOS=linux GOARCH=386 go build
scp gh-docker dev:gh-docker
ssh $SSH_HOST sudo mv gh-docker /bin/gh-docker
