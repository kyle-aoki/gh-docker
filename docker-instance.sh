###################################################################################################
#### files ########################################################################################
###################################################################################################
curl https://raw.githubusercontent.com/kyle-aoki/petra/main/petra             > petra
curl https://raw.githubusercontent.com/kyle-aoki/petra/main/petra.service     > petra.service
curl https://raw.githubusercontent.com/kyle-aoki/petra/main/petra-config.json > petra-config.json
curl https://raw.githubusercontent.com/kyle-aoki/petra/main/nginx.conf        > nginx.conf
curl https://raw.githubusercontent.com/kyle-aoki/petra/main/app-config.json   > app-config.json

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

sudo mkdir /petra
sudo mv petra-config.json /petra/petra-config.json
sudo chmod 600 /petra/petra-config.json
sudo mv app-config.json /petra/app-config.json
sudo chmod 600 /petra/app-config.json

###################################################################################################
#### certbot ######################################################################################
###################################################################################################
sudo apt install certbot -y
sudo apt install python3-certbot-nginx -y
sudo certbot --nginx --agree-tos --non-interactive -m example@example.com --domain subdomain.domain.com

###################################################################################################
#### systemctl ####################################################################################
###################################################################################################
sudo mv petra.service /etc/systemd/system/petra.service
sudo systemctl daemon-reload
sudo systemctl enable petra
sudo systemctl start petra
sudo systemctl status petra
