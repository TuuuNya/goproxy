package provider

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

type GithubSocks5 struct{}

func (g *GithubSocks5) Name() string {
	return "github.com"
}

func (g *GithubSocks5) SupportedTypes() []string {
	return []string{"socks5"}
}

func (g *GithubSocks5) FetchProxies(types []string) (proxies []Proxy, err error) {
	if !CheckTypeSupported(g.SupportedTypes(), types) {
		return nil, nil
	}

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		// every line is a proxy in the format: IP:PORT
		proxyStrs := strings.Split(string(r.Body), "\r\n")
		for _, proxyStr := range proxyStrs {
			// Split IP and port
			parts := strings.Split(proxyStr, ":")
			if len(parts) != 2 {
				continue
			}

			// Check if the proxy type is supported
			for _, t := range types {
				if t == "socks5" {
					proxies = append(proxies, Proxy{
						IP:   parts[0],
						Port: parts[1],
						Type: "socks5",
					})
				}
			}
		}
	})

	c.Visit("https://raw.githubusercontent.com/hookzof/socks5_list/refs/heads/master/proxy.txt")

	return proxies, nil
}
