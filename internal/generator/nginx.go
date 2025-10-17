package generator

import (
	"fmt"
	"ldriko/dokploy-bob/internal/config"
	"ldriko/dokploy-bob/internal/exporter"
)

type NginxConfig struct {
	Target   string `yaml:"-"`
	Services map[string]NginxService
}

type NginxService struct {
	ServerName []string
	PHP        config.PHPConfig
	AccessLog  string
	ErrorLog   string
}

func NewNginxConfig() *NginxConfig {
	return &NginxConfig{
		Services: make(map[string]NginxService),
	}
}

func (nc *NginxConfig) SetTarget(target string) {
	nc.Target = target
}

func (nc *NginxConfig) AddService(name string, svc *config.Service, pi config.ProviderInstance) error {
	service := NginxService{
		ServerName: svc.Domains,
		AccessLog:  fmt.Sprintf("/var/log/nginx/%s.access.log", name),
		ErrorLog:   fmt.Sprintf("/var/log/nginx/%s.error.log", name),
	}

	if npc, ok := pi.Config.(*config.NginxProviderConfig); ok {
		service.PHP.Root = npc.Root
		if npc.Type == "php" {
			service.PHP.Version = npc.PHP.Version
		}
	}

	nc.Services[name] = service

	return nil
}

func (nc *NginxConfig) Export(path string) error {
	for name, service := range nc.Services {
		filename := name + ".conf"
		serverNames := ""
		for i, domain := range service.ServerName {
			if i > 0 {
				serverNames += " "
			}
			serverNames += domain
		}

		var data string
		if service.PHP.Version != "" {
			data = fmt.Sprintf(`server {
    listen 8080;
    server_name %s;

    root %s;
    index index.php index.html index.htm;

    access_log %s;
    error_log %s;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \\.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/run/php/php%s-fpm.sock;

        # Forward proxy headers to PHP-FPM
        fastcgi_param HTTP_X_FORWARDED_FOR $proxy_add_x_forwarded_for;
        fastcgi_param HTTP_X_FORWARDED_PROTO $http_x_forwarded_proto;
        fastcgi_param HTTP_X_FORWARDED_HOST $http_x_forwarded_host;
        fastcgi_param HTTP_X_FORWARDED_PORT $http_x_forwarded_port;

        # Mark HTTPS if proxied as https
        set $https_off "";
        if ($http_x_forwarded_proto = "https") {
            set $https_off on;
        }
        fastcgi_param HTTPS $https_off;
        
        include fastcgi_params;
    }

    location ~ /\\.ht {
        deny all;
    }
}
`,
				serverNames,
				service.PHP.Root,
				service.AccessLog,
				service.ErrorLog,
				service.PHP.Version,
			)
		} else {
			data = fmt.Sprintf(`server {
    listen 8080;
    server_name %s;

    root %s;
    index index.html index.htm;

    access_log %s;
    error_log %s;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location ~ /\\.ht {
        deny all;
    }
		
	location ~* \.(?:js|css|png|jpg|jpeg|gif|ico|svg|woff2?)$ {
		expires 1y;
		add_header Cache-Control "public, immutable";
	}
}
`,
				serverNames,
				service.PHP.Root,
				service.AccessLog,
				service.ErrorLog,
			)
		}

		err := exporter.Process(path+"/"+filename, []byte(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func (nc *NginxConfig) GetTarget() string {
	return nc.Target
}
