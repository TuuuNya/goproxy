package provider

import "strings"

type Proxy struct {
	IP   string
	Port string
	Type string
}

type Provider interface {
	Name() string
	SupportedTypes() []string
	FetchProxies(types []string) ([]Proxy, error)
}

// Check if the proxy type is supported
func CheckTypeSupported(types []string, specified_types []string) bool {
	for _, t := range types {
		for _, st := range specified_types {
			if t == strings.ToLower(st) {
				return true
			}
		}
	}
	return false
}
