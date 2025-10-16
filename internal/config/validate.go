package config

import "fmt"

func (c *Config) Validate() error {
	if len(c.Providers) == 0 {
		return fmt.Errorf("at least one target must be specified")
	}
	if len(c.Services) == 0 {
		return fmt.Errorf("at least one site must be specified")
	}
	for name, serv := range c.Services {
		if len(serv.Domains) == 0 {
			return fmt.Errorf("site %s must have at least one domain", name)
		}

		for _, pi := range serv.Providers {
			valid := false
			for _, v := range providerTypeNames {
				if pi.Name == v {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("site %s has an invalid provider: %s", name, pi.Name)
			}
			if pi.Config != nil {
				if err := pi.Config.Validate(); err != nil {
					return fmt.Errorf("site %s %s provider: %v", name, pi.Name, err)
				}
			}
		}
	}
	return nil
}
