# Pi Temp

[![Build Status](https://travis-ci.org/andrewmarklloyd/pi-temp.svg?branch=master)](https://travis-ci.org/andrewmarklloyd/pi-temp)

Small web server running on a Raspberry Pi that get's temperature readings.


### One Line Install
To install on a Raspberry Pi with a single line command, run the following:
```
bash <(curl -s -H 'Cache-Control: no-cache' https://raw.githubusercontent.com/andrewmarklloyd/pi-temp/master/install/install.sh)
```

### Developing Locally
Requires Go 1.13.1 to build the project.
```
# run the program
go run main.go

# build an executable
go build -o pi-temp main.go
```
