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
			// trim socks5:// prefix
			proxyStr = strings.TrimPrefix(proxyStr, "socks5://")
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
	c.Visit("https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/socks5/data.txt")
	c.Visit("https://raw.githubusercontent.com/roosterkid/openproxylist/refs/heads/main/SOCKS5_RAW.txt")
	c.Visit("https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/refs/heads/master/socks5.txt")
	c.Visit("https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/refs/heads/main/socks5_checked.txt")
	c.Visit("https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-socks5.txt")
	c.Visit("https://raw.githubusercontent.com/r00tee/Proxy-List/refs/heads/main/Socks5.txt")
	c.Visit("https://raw.githubusercontent.com/officialputuid/KangProxy/refs/heads/KangProxy/socks5/socks5.txt")
	c.Visit("https://raw.githubusercontent.com/dpangestuw/Free-Proxy/refs/heads/main/socks5_proxies.txt")
	c.Visit("https://raw.githubusercontent.com/BreakingTechFr/Proxy_Free/refs/heads/main/proxies/socks5.txt")
	c.Visit("https://raw.githubusercontent.com/fyvri/fresh-proxy-list/archive/storage/classic/socks5.txt")
	c.Visit("https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/refs/heads/main/socks5.txt")
	c.Visit("https://raw.githubusercontent.com/yemixzy/proxy-list/main/proxies/socks5.txt")

	return proxies, nil
}
