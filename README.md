# MQTT client V-0.1

Broker [eclipse mosquitto](https://hub.docker.com/_/eclipse-mosquitto):

``` bash
docker run -it -p 1883:1883 -p 9001:9001 -v /run/mosquitto/data:/mosquitto/data -v /mosquitto/log eclipse-mosquitto
```
