package provider

import (
	"encoding/base64"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type FreeProxyCz struct{}

func (f *FreeProxyCz) Name() string {
	return "free-proxy.cz"
}

func (f *FreeProxyCz) SupportedTypes() []string {
	return []string{"http", "https", "socks4", "socks5"}
}

func (f *FreeProxyCz) FetchProxies(types []string) (proxies []Proxy, err error) {
	if !CheckTypeSupported(f.SupportedTypes(), types) {
		return nil, nil
	}

	c := colly.NewCollector(
		colly.AllowedDomains("free-proxy.cz"),
		colly.URLFilters(
			regexp.MustCompile(`http://free-proxy.cz/en/proxylist/.*`),
		),
	)

	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach("#proxy_list tbody tr", func(_ int, el *colly.HTMLElement) {
			tds := el.ChildTexts("td")
			if len(tds) == 11 {
				port := tds[1]
				proxyType := strings.ToLower(tds[2])

				encodedIp := tds[0]
				encodedIp = strings.ReplaceAll(encodedIp, "document.write(Base64.decode(\"", "")
				encodedIp = strings.ReplaceAll(encodedIp, "\"))", "")
				ip, err := base64.StdEncoding.DecodeString(encodedIp)
				if err == nil {
					for _, t := range types {
						if t == proxyType {
							proxies = append(proxies, Proxy{
								IP:   string(ip),
								Port: port,
								Type: proxyType,
							})
						}
					}
				}
			}
		})

		e.ForEach(".paginator a", func(_ int, el *colly.HTMLElement) {
			c.Visit(el.Request.AbsoluteURL(el.Attr("href")))
		})

	})

	c.Visit("http://free-proxy.cz/en/proxylist/country/all/all/ping/all")

	return proxies, nil
}
