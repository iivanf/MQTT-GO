# MQTT client V-0.1

### Envio de datos dende barco e lectura dende estación

Para desplegar o broker mqtt usamos o propio docker proporcionado por eclipse-mosquitto, dandolle un volumen a data para obter persistencia nas mensaxes, redirixindo, asimesmo, os portos internos do docker os mesmos no host:

``` bash
docker run -it -p 1883:1883 -p 9001:9001 -v /run/mosquitto/data:/mosquitto/data -v /mosquitto/log eclipse-mosquitto
```