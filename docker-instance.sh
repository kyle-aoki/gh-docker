###################################################################################################
#### files ########################################################################################
###################################################################################################
curl https://raw.githubusercontent.com/kyle-aoki/petra/main/petra              > petra
curl https://raw.githubusercontent.com/kyle-aoki/petra/main/petra-config.json  > petra-config.json
curl https://raw.githubusercontent.com/kyle-aoki/petra/main/nginx.conf         > nginx.conf
curl https://raw.githubusercontent.com/kyle-aoki/petra/main/config.json        > config.json

###################################################################################################
#### install ######################################################################################
###################################################################################################
sudo apt update -y
sudo apt install docker.io -y
sudo chmod 666 /var/run/docker.sock

sudo apt install nginx -y
sudo mv nginx.conf /etc/nginx/nginx.conf
sudo nginx -s reload

sudo chmod 700 petra
sudo mv petra /bin/petra

###################################################################################################
#### https ########################################################################################
###################################################################################################
sudo apt install certbot -y
sudo apt install python3-certbot-nginx -y
sudo certbot --nginx --agree-tos --non-interactive -m example@example.com --domains subdomain.domain.com
