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

func (nc *NginxConfig) AddService(name string, svc *config.Service) error {
	service := NginxService{
		ServerName: svc.Domains,
		AccessLog:  fmt.Sprintf("/var/log/nginx/%s.access.log", name),
		ErrorLog:   fmt.Sprintf("/var/log/nginx/%s.error.log", name),
	}

	service.PHP.Version = svc.PHP.Version
	service.PHP.Root = svc.PHP.Root

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

		data := fmt.Sprintf(`server {
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
