package provider

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

type ProxyListPlus struct{}

func (p *ProxyListPlus) Name() string {
	return "list.proxylistplus.com"
}

func (p *ProxyListPlus) SupportedTypes() []string {
	return []string{"http", "https", "socks5"}
}

func (p *ProxyListPlus) FetchProxies(types []string) (proxies []Proxy, err error) {
	if !CheckTypeSupported(p.SupportedTypes(), types) {
		return nil, nil
	}

	c := colly.NewCollector(
		colly.AllowedDomains("list.proxylistplus.com"),
	)

	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach("table tr.cells", func(_ int, e *colly.HTMLElement) {
			tds := e.ChildTexts("td")

			if len(tds) == 8 {
				IP := tds[1]
				Port := tds[2]
				Type := strings.ToLower(tds[3])
				Https := strings.ToLower(tds[6])

				if Type != "socks5" {
					if Https == "yes" {
						Type = "https"
					} else {
						Type = "http"
					}
				} else {
					Type = "socks5"
				}

				for _, t := range types {
					if t == Type {
						proxies = append(proxies, Proxy{
							IP:   IP,
							Port: Port,
							Type: Type,
						})
					}
				}
			}

		})
	})

	c.Visit("https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-1")
	c.Visit("https://list.proxylistplus.com/Socks-List-1")
	c.Visit("https://list.proxylistplus.com/SSL-List-1")

	return proxies, nil
}
