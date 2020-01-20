#!/bin/bash

archive_path="/tmp/pi-temp"
mkdir -p ${archive_path}

latestVersion=$(curl --silent "https://api.github.com/repos/andrewmarklloyd/pi-temp/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
curl -sL https://github.com/andrewmarklloyd/pi-temp/archive/${latestVersion}.tar.gz | tar xvfz - -C "${archive_path}" --strip 1 > /dev/null

binaryUrl=$(curl -s https://api.github.com/repos/andrewmarklloyd/pi-temp/releases/latest | jq -r '.assets[] | select(.name == "pi-temp") | .browser_download_url')
curl -sL $binaryUrl -o ${archive_path}/pi-temp
chmod +x ${archive_path}/pi-temp
rm -f ./install/*
rm -f ./static/*
cp ${archive_path}/install/* install/
cp ${archive_path}/static/* static/
echo -n ${latestVersion} > /home/pi/static/version
echo -n ${latestVersion} > /home/pi/static/latestVersion
mv ${archive_path}/pi-temp ./
rm -rf ${archive_path}
sudo systemctl restart pi-temp.service
