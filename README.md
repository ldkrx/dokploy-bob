# Dokploy Bob

Bob will generate traefik and nginx config from a yaml file like this:

```yaml
providers:
  traefik: 
    target: ./tests/results/traefik.yaml
  nginx: 
    target: ./tests/results/nginx

services:
  service-1:
    domains:
      - example.com
      - www.example.com
    providers:
      - traefik
      - nginx
    port: 8080
    php:
      version: 8.1
      root: /var/www/service-1/public

  service-2:
    domains:
      - blog.example.com
      - www.blog.example.com
    providers: 
      - traefik
      - node
    port: 3000
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
  services:
    service-1:
      loadBalancer:
        servers:
          - url: http://service-1:8080
    service-2:
      loadBalancer:
        servers:
          - url: http://service-2:3000
```

**`/etc/nginx/sites-available/generated/service-1.conf`**:

```
server {
    listen 8080;
    server_name example.com www.example.com;

    root /public;
    index index.php index.html index.htm;

    access_log /var/log/nginx/service-1.access.log;
    error_log /var/log/nginx/service-1.error.log;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \\.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/run/php/8.1-fpm.sock;
    }

    location ~ /\\.ht {
        deny all;
    }
}
```

#### When is Bob useful:

- Need to serve Nginx (:8080), Node (:4xxx), and Dokploy (:80 and :443)
- Too lazy to setup these configs
- Using Traefik as a reverse proxy to Nginx / Node / basically anything