## go-proxy

### How to use it? 

1. rename go-proxy.yaml.example -> go-proxy.yaml
2. update the targets
3. and docker run .

### docker compose

```
version: '3'
services:
	  x-proxy:
    container_name: x-proxy
    hostname: x-proxy
    build:
      context: ./
      dockerfile: Dockerfile
    ports:
      - '1988:1988'
```
