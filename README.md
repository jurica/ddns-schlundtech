# ddns-schlundtech
Simple self hosted ddns service for [SchlundTech](https://www.schlundtech.com) users, making use of the [XML-Gateway](https://www.schlundtech.com/services/xml-gateway/).

I've created this to suite _my_ needs, so it's far from being "feature complete" and some assumtions might break usage for others.

Feedback and contributions are always welcome!

## What it does
1. Provides an _unsecured_ http endpoint wating for update request.
1. On request:
    1. Zone Info (0205) is requested for the given domain.
        1. The given domain is split at the first "." into two parts for this.
        1. First part is interpreted as ressource record name (aka "subdomain").
        1. Second part is interpreted as "zone".
    1. Zone Update (Bulk) 0202001 is sent that
        1. deletes the existing ressource record
        1. creates a new one with the given IP

Currently there is no advances checks or logic etc. So it's initially required to setup the subdomain in the SchlundTech panel first.

## Setup
### Build
```
docker build -t ddns-schlundtech .
```
This will create an alpine based container with the binary.
### Install
See [docker-compose.yml](docker-compose.yml) for how to run the container.

Prepare a save location for the configuration file with the following content:
```
user = ""
password = ""
context = ""
```
This file needs to be bind mounted into the container, providing the necessary runtime configuration.

### OpenWRT
I've just tested with [OpenWRT](https://openwrt.org). The following settings are working for me:

| Section ||
| --- | --- |
| **Basic Settings** ||
| DDNS Service Provider | -- custom -- |
| Lookup Hostname | {DOMAIN} |
| IP address version | IPv4-Address |
| Custom update-URL | {DOCKER_HOST}:8080?ip=[IP]&domain=[DOMAIN] |
| Custom update-script | `empty` |
| Domain | {DOMAIN} |
| Username | none |
| Password | none |
| Optional Encoded Parameter | `empty` |
| Optional Parameter | `empty` |
| Use HTTP Secure | `unchecked` |
| **Advances Settings** ||
| IP address source | Interface |
| Interface | pppoe-WAN |

Replace placeholders: {DOMAIN}, {DOCKER_HOST}