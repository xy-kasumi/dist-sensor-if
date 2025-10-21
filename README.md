# dist-sensor-if

Interfaces to Optex FA CD22-15-485M12 sensor via RasPi 5.
The program provides phone-friendly web interface.

## Prerequisite
WSL2 (host)
* mDNS is working
* podman and Go is installed
* can ssh to host with keys (no password)

RasPi 5 (device)
* Hostname is `dist-sensor.local`
* Podman is installed

## Deployment
* Run `./build.sh` to generate container image
* Run `./deploy.sh` to reload the container image in the device
