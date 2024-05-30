###################################################################################################
#### files ########################################################################################
###################################################################################################
SSH_HOST=dev

GOOS=linux GOARCH=386 go build

scp petra              "$SSH_HOST":petra
scp petra-config.json  "$SSH_HOST":petra-config.json
scp nginx.conf         "$SSH_HOST":nginx.conf
scp config.json        "$SSH_HOST":config.json

###################################################################################################
#### install ######################################################################################
###################################################################################################
sudo apt update -y
sudo apt install docker.io -y
sudo chmod 666 /var/run/docker.sock

sudo apt install nginx -y
sudo mv nginx.conf /etc/nginx/nginx.conf
sudo nginx -s reload

sudo mv petra /bin/petra

###################################################################################################
#### https ########################################################################################
###################################################################################################
sudo apt install certbot -y
sudo apt install python3-certbot-nginx -y
sudo certbot --nginx --agree-tos --non-interactive -m kyle.aoki.000@gmail.com --domains student-api-dev.kyle-aoki.dev
