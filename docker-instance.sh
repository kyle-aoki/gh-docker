SSH_HOST=dev

function run {
  ssh $SSH_HOST "$@"
}

###################################################################################################
#### dependencies #################################################################################
###################################################################################################
run sudo apt update -y
run sudo apt install docker.io -y
run sudo chmod 666 /var/run/docker.sock
run sudo apt install nginx -y

###################################################################################################
#### install ######################################################################################
###################################################################################################
scp ./nginx.conf $SSH_HOST:nginx.conf
run sudo mv nginx.conf /etc/nginx/nginx.conf
run sudo nginx -s reload

GOOS=linux GOARCH=386 go build
scp petra $SSH_HOST:petra
scp petra-config.json $SSH_HOST:petra-config.json
run sudo mv petra /bin/petra
scp config.json "$SSH_HOST":config.json
