###################################################################################################
#### nginx ########################################################################################
###################################################################################################

# if "docker ps --format=json | jq .Ports" has 4 "8080"s, then we are currently on port 8080
function nginx_next_proxy_port {
  port_8080_count=$(docker ps --format=json | jq .Ports | grep -Fo 8080 | wc -l)

  if [ $port_8080_count == "0" ]; then
    echo -n "8080"
  elif [ $port_8080_count == "2" ]; then
    echo -n "8080"
  elif [ $port_8080_count == "4" ]; then
    echo -n "8081"
  else
    echo "unknown"
  fi
}

###################################################################################################
#### deployment ###################################################################################
###################################################################################################

function deploy {
  local docker_username=$1
  local docker_password=$2
  local docker_image_repository=$3
  local docker_image_tag=$4
  local docker_tag=$docker_username/$docker_image_repository:$docker_image_tag

  echo "logging into docker"
  docker login --username=$docker_username --password=$docker_password

  echo "pulling tag: $docker_tag"
  docker pull $docker_tag

  docker_current_ps=$(docker ps -q)
  echo "found existing container ids: $docker_current_ps"

  next_port=$(nginx_next_proxy_port)
  echo "found next port: $next_port"

  if [ $next_port == "unknown" ]; then
    echo "unknown port state"
    return 1
  fi

  echo "attempting to run new image tag"
  docker run -d -p $next_port:8080 --restart=always $docker_tag || exit 1

  echo "switching nginx config and reloading"
  sudo cp nginx-$next_port.conf /etc/nginx/nginx.conf
  sudo nginx -s reload

  echo "attempting to stop and terminate old containers"
  if [ "$docker_current_ps" ]; then
    docker stop $docker_current_ps
    docker rm $docker_current_ps
  fi
}

# deploy kyleaoki500 xxx gh-docker BRAVO-5
# deploy ${{ vars.docker_username }} ${{ secrets.docker_password }} ${{ vars.docker_image_repository }} ${{ vars.docker_image_tag }}
