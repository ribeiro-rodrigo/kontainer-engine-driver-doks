# Kontainer Engine DOKS (DigitalOcean Kubernetes Service) Driver
[![Build Status](https://travis-ci.org/ribeiro-rodrigo/kontainer-engine-driver-doks.svg?branch=master)](https://travis-ci.org/ribeiro-rodrigo/kontainer-engine-driver-doks)
[![godoc](https://godoc.org/github.com/ribeiro-rodrigo/kontainer-engine-driver-doks?status.svg)](https://godoc.org/github.com/ribeiro-rodrigo/kontainer-engine-driver-doks)
[![Coverage](https://codecov.io/gh/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/branch/master/graph/badge.svg)](https://codecov.io/gh/ribeiro-rodrigo/kontainer-engine-driver-digitalocean)

This repo contains the DOKS (DigitalOcean Kubernetes Service) driver for the rancher server

## Building
```shell script
make
```
Will output driver binaries into the dist directory, these can be imported directly into Rancher and used as cluster drivers. They must be distributed via URLs that your Rancher instance can establish a connection to and download the driver binaries.

## Running Local
```shell script
./dist/kontainer-engine-driver-digitalocean-linux $PORT
```
or
```shell script
./dist/kontainer-engine-driver-digitalocean-darwin $PORT
```

## Installing in Rancher
Go to the Cluster Drivers management screen in Rancher and click Add Cluster Driver. Enter the URL of your driver, a UI URL (see the [UI repo](https://github.com/ribeiro-rodrigo/ui-cluster-driver-doks) for details), and a checksum (optional), and click Create. Rancher will automatically download and install your driver. It will then become available to use on the Add Cluster screen.