package provider

import (
	"encoding/base64"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type ProxyListOrg struct{}

func (p *ProxyListOrg) Name() string {
	return "proxy-list.org"
}

func (p *ProxyListOrg) SupportedTypes() []string {
	return []string{"http", "https"}
}

func (p *ProxyListOrg) FetchProxies(types []string) (proxies []Proxy, err error) {
	// Check if the proxy type is supported
	if !CheckTypeSupported(p.SupportedTypes(), types) {
		return nil, nil
	}

	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(`^https://proxy-list\.org/english/index\.php\?p=\d+$`),
		),
	)

	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach(".proxy-table .table ul", func(i int, e *colly.HTMLElement) {
			proxyStr := e.ChildText(".proxy")
			if proxyStr == "" {
				return
			}

			decodedProxy, err := decodeProxy(proxyStr)
			if err != nil {
				return
			}

			// Split IP and port
			parts := strings.Split(decodedProxy, ":")
			if len(parts) != 2 {
				return
			}

			// Check if the proxy type is supported
			proxyType := strings.ToLower(e.ChildText(".https"))
			for _, t := range types {
				if t == proxyType {
					proxies = append(proxies, Proxy{
						IP:   parts[0],
						Port: parts[1],
						Type: proxyType,
					})
				}
			}
		})

		e.ForEach(".table-menu a[href]", func(_ int, e *colly.HTMLElement) {
			c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		})
	})

	c.Visit("https://proxy-list.org/english/index.php?p=1")

	return proxies, nil
}

func decodeProxy(proxyStr string) (string, error) {
	cleanedStr := strings.TrimPrefix(proxyStr, "Proxy('")
	cleanedStr = strings.TrimSuffix(cleanedStr, "')")

	// Base64 解码
	decodedBytes, err := base64.StdEncoding.DecodeString(cleanedStr)
	if err != nil {
		return "", err
	}

	return string(decodedBytes), nil
}
