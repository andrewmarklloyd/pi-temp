# Pi Temp
Small web server running on a Raspberry Pi that monitors temperature.

### Build and Deploy
Requires Go 1.13.1 to build the project. The `util.sh` script will build and copy the binary to a remote Raspberry Pi.
```
# build an executable
./util.sh build

# copy executable and systemd config, start service
./util.sh pi raspberrypi.local
```
