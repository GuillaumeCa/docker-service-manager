# Docker service manager

Web ui to manage Docker services

## Features
- View logs in real time
- Update all services (with new images, support private registry)
  - Fetch new image (private registry supported)
  - Service blacklist can be configured in the config.json file

## Launch image with docker-compose

```yaml
version: "3"
services:
  docker-manager:
    image: docker-manager
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /config/config.json:/config/config.json
```

Default config file :
```json
{
  "auth": {
    "user": "admin",
    "password": "password"
  },
  "registry": {
    "url": "localhost",
    "user": "admin",
    "password": "password"
  },
  "blacklist": []
}
```