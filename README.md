# Dokploy Bob

Bob will generate traefik and nginx config from a yaml file like this:

```yaml
providers:
  traefik:
    target: /etc/traefik/dynamics/host.yaml
  nginx:
    target: /etc/nginx/sites-available/generated
  node:
    target: /root/node-apps.json

services:
  service-1:
    domains:
      - example.com
      - www.example.com
    port: 8080
    providers:
      - traefik
      - nginx:
          type: php
          root: /var/www/service-1/public
          php:
            version: 8.1

  service-2:
    domains:
      - blog.example.com
      - www.blog.example.com
    port: 3001
    providers:
      - traefik
      - node:
          script: /var/www/service-2/server.js
          interpreter: /root/.nvm/versions/node/v22.20.0/bin/node
          post-update:
            - npm install

  service-3:
    domains:
      - static.example.com
    port: 8080
    providers:
      - traefik
      - nginx:
          type: static
          root: /var/www/service-3/dist
```

To This:

**`/etc/dokploy/traefik/dynamic/host.yaml`**:

```yaml
http:
  routers:
    service-1:
      rule: Host(`example.com`) || Host(`www.example.com`)
      entryPoints:
        - web
        - websecure
      service: service-1
      tls:
        certResolver: letsencrypt
    service-2:
      rule: Host(`blog.example.com`) || Host(`www.blog.example.com`)
      entryPoints:
        - web
        - websecure
      service: service-2
      tls:
        certResolver: letsencrypt
    service-3:
      rule: Host(`static.example.com`)
      entryPoints:
        - web
        - websecure
      service: service-3
      tls:
        certResolver: letsencrypt
  services:
    service-1:
      loadBalancer:
        servers:
          - url: http://172.17.0.1:8080
    service-2:
      loadBalancer:
        servers:
          - url: http://172.17.0.1:3001
    service-3:
      loadBalancer:
        servers:
          - url: http://172.17.0.1:8080
```

**`/etc/nginx/sites-available/generated/service-1.conf`**:

```conf
server {
    listen 8080;
    server_name example.com www.example.com;

    root /var/www/service-1/public;
    index index.php index.html index.htm;

    access_log /var/log/nginx/service-1.access.log;
    error_log /var/log/nginx/service-1.error.log;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \\.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/run/php/php8.1-fpm.sock;
    }

    location ~ /\\.ht {
        deny all;
    }
}
```

**`/etc/nginx/sites-available/generated/service-3.conf`**:

```conf
server {
    listen 8080;
    server_name static.example.com;

    root /var/www/service-3/dist;
    index index.html index.htm;

    access_log /var/log/nginx/service-3.access.log;
    error_log /var/log/nginx/service-3.error.log;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location ~ /\\.ht {
        deny all;
    }
}
```

**`/root/node-apps.json`**

```json
{
    "apps": [
        {
            "name": "service-2",
            "script": "/var/www/service-2/server.js",
            "interpreter": "/root/.nvm/versions/node/v22.20.0/bin/node",
            "env": {
                "PORT": 3001
            }
        }
    ]
}
```

#### When is Bob useful:

- Need to serve Nginx (PHP or static sites on :8080), Node (:3xxx), and Dokploy (:80 and :443)
- Too lazy to setup these configs
- Using Traefik as a reverse proxy to Nginx / Node / basically anything

#### Nginx Types:
- `php`: For PHP applications with FastCGI
- `static`: For static websites (HTML, JS, CSS)
