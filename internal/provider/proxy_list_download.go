package provider

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

type ProxyListDownload struct{}

func (p *ProxyListDownload) Name() string {
	return "proxy-list.download"
}

func (p *ProxyListDownload) SupportedTypes() []string {
	return []string{"http", "https", "socks5"}
}

func (p *ProxyListDownload) FetchProxies(types []string) (proxies []Proxy, err error) {
	if !CheckTypeSupported(p.SupportedTypes(), types) {
		return nil, nil
	}

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		// URL query parameter is the proxy type
		currentType := r.Request.URL.Query().Get("type")

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
				if t == currentType {
					proxies = append(proxies, Proxy{
						IP:   parts[0],
						Port: parts[1],
						Type: currentType,
					})
				}
			}
		}
	})

	for _, t := range types {
		c.Visit(fmt.Sprintf("https://www.proxy-list.download/api/v1/get?type=%s", t))
	}

	return proxies, nil
}
