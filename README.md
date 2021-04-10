# PAN-OS UserID transaction monitor
Palo Alto Networks NGFW's feature a very powerfull API endpoint called User-ID.

It allows tagging IP address objects as well as User objects. These tags can then
be used to dynamically group IP addresses (Dynamic Address Group - DAG) as well as to
dynamically group users (Dynamic User Group - DUG). Resulting groups can be consumed
in several policies like session security or packet forwarding.

The User-ID API allows, as well, registering these user objects (User Login/Logout) and
these IP address objects (Object Register/Unregister)

One of the most powerful details of the corresponding PAN-OS feature is both registration
and tag transactions can contain a timeout value. This way the user can, for instance,
register a user login with a 10 minute expiration mark and then leave the NGFW take care of
removing the entry at its due time.

There are tons of examples on how to script the User-ID API so these groups could be
created from SNMP alerts, SYSLOG events or application log messages.

As the use case grows (in number to objects and tags) the need for a easy way to query
the current state (of group membership) raises in urgency. Although PAN-OS provides CLI, WebUI
and XML API ways of displaying the current state I always missed a basic
tag-to-ip API endpoint that could dump the corresponding list in plain text format.

This small "man-in-the-midle" proxy application monitors User-ID API payloads and
keeps its own copy of registered objects and tags. It tracks due date of the objects and
tags using its own internal clock reference (unsynchonised from the NGFW).

Then it exposes the basic dump endpoint in `GET /edl/?list=[user|group|tag]&key=<tag>`

## Running the application
There is only one mandatory `TARGET` environmental variable configuration value to be
provided in the form of `host:port` (FQDN or Address of the PAN-OS device to proxy)

```sh
TARGET="<ip:port>" go run .
```

A http service would be started in the port `8080` of localhost. The http service will
forward all messages to the NGFW (including the User-ID API ones). But it will monitor
User-ID transactions keeping its internal state alive.

Use the `http://127.0.0.1:8080/edl/?list=[user|group|tag]&key=<tag>` URL to dump the 
state of the different lists.

The application would honour a different TCP port as provided in the `PORT` environmental
variable. In case the NGFW device is using a self-signed or untrusted certificate
then you can configure the proxy to ignore the certificate chain with `INSECURE=true`

```sh
TARGET="<ip:port>" PORT=8081 INSECURE=true go run .
```

## Microservice version
A Dockerfile is provided to pack the code in a distroless container.

Or you can use the prebuild docker image

```sh
docker pull ghcr.io/xhoms/uidmonitor:latest
```