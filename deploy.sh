###################################################################################################
#### deploy #######################################################################################
###################################################################################################

# if "docker ps --format=json | jq .Ports" has 4 "8080"s, then we are currently on port 8080
function find_next_proxy_port {
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

function deploy {
  local docker_username=$1
  local docker_password=$2
  local docker_image_repository=$3
  local docker_image_tag=$4

  local docker_tag=$docker_username/$docker_image_repository:$docker_image_tag

  docker login --username=$docker_username --password=$docker_password
  docker pull $docker_tag

  docker_current_ps=$(docker ps -q)
  echo "found existing container ids: $docker_current_ps"

  next_port=$(find_next_proxy_port)
  echo "found next port: $next_port"

  if [ $next_port == "unknown" ]; then
    echo "unknown port state"
    return 1
  fi

  docker run -d -p $next_port:8080 --restart=always $docker_tag || exit 1

  if [ "$docker_current_ps" ]; then
    echo "attempting to stop and terminate old containers"
    docker stop $docker_current_ps
    docker rm $docker_current_ps
  else
    echo "no old containers found"
  fi
}

deploy $1 $2 $3 $4
