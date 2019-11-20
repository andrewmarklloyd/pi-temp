#!/bin/bash

build() {
  GOOS=linux GOARCH=arm GOARM=5 go build -o pi-temp main.go
}

install() {
  scp -pr ./static ${username}@${host}:
  scp pi-temp ${username}@${host}:
  scp pi-temp.service ${username}@${host}:
  ssh ${username}@${host} "sudo mv pi-temp.service /etc/systemd/system/; sudo systemctl enable pi-temp.service; sudo systemctl start pi-temp.service"
}

deploy() {
  ssh ${username}@${host} "sudo systemctl stop pi-temp.service"
  scp pi-temp ${username}@${host}:
  scp -pr ./static ${username}@${host}:
  ssh ${username}@${host} "sudo systemctl restart pi-temp.service"
}

if [[ ${1} == 'install' ]]; then
  build
  install
elif [[ ${1} == 'build' ]]; then
  build
elif [[ ${1} == 'deploy' ]]; then
  username=${2}
  host=${3}
  if [[ -z ${username} || -z ${host} ]]; then
    echo "Use username and host as arguments. Example:"
    echo "./util.sh deploy pi raspberrypi.local"
    exit 1
  fi
  build
  deploy
fi
