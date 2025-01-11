package provider

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type XseoIn struct{}

func (x *XseoIn) Name() string {
	return "xseo.in"
}

func (x *XseoIn) SupportedTypes() []string {
	return []string{"http", "socks5"}
}

func (x *XseoIn) FetchProxies(types []string) (proxies []Proxy, err error) {
	if !CheckTypeSupported(x.SupportedTypes(), types) {
		return nil, nil
	}

	c := colly.NewCollector(
		colly.AllowedDomains("xseo.in"),
	)

	// Handle the HTML response
	c.OnHTML("html", func(e *colly.HTMLElement) {
		js_value_map := make(map[string]string)

		e.ForEach("script", func(i int, t *colly.HTMLElement) {
			if regexp.MustCompile(`^(\w=\d;){10}$`).MatchString(t.Text) {
				js_value_map = parseJSValues(t.Text)
				if len(js_value_map) == 0 {
					return
				}

				e.ForEach("tr", func(i int, t *colly.HTMLElement) {
					// if t has 6 children, it's a proxy row
					if len(t.ChildTexts("td")) == 6 {
						proxySlice := t.ChildTexts("td")
						proxyStr := proxySlice[0]

						if regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+:.*$`).MatchString(proxyStr) {
							decodeProxyStr := decodeProxyStr(proxyStr, js_value_map)

							ip := strings.Split(decodeProxyStr, ":")[0]
							port := strings.Split(decodeProxyStr, ":")[1]
							proxyType := strings.ToLower(proxySlice[2])

							for _, t := range types {
								if t == proxyType {
									proxies = append(proxies, Proxy{
										IP:   ip,
										Port: port,
										Type: proxyType,
									})
								}
							}

						}
					}
				})

			}
		})
	})

	// Perform the POST request manually
	postData := []byte("submit=1")
	headers := http.Header{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}

	err = c.Request("POST", "https://xseo.in/proxylist", bytes.NewReader(postData), nil, headers)
	if err != nil {
		return nil, err
	}

	return proxies, nil
}

func parseJSValues(js string) map[string]string {
	js_value_map := make(map[string]string)
	js_value_re := regexp.MustCompile(`^(\w)=(\d)$`)

	for _, line := range bytes.Split([]byte(js), []byte(";")) {
		if js_value_re.FindAllSubmatch(line, -1) != nil {
			matches := js_value_re.FindAllSubmatch(line, -1)
			js_value_map[string(matches[0][1])] = string(matches[0][2])
		}
	}

	return js_value_map
}

func decodeProxyStr(encoded string, js_values map[string]string) string {
	// 93.157.248.108:document.write(""+j+j)
	proxySlice := strings.Split(encoded, ":")

	ip := proxySlice[0]

	encodedPort := proxySlice[1]
	encodedPort = strings.ReplaceAll(encodedPort, "document.write(", "")
	encodedPort = strings.ReplaceAll(encodedPort, ")", "")

	// port是js加密了，document.write(""+j+j)这种形式
	// 解密
	decodedPort := ""
	for _, c := range encodedPort {
		if c == '+' || c == '"' {
			continue
		}
		decodedPort = fmt.Sprintf("%s%s", decodedPort, js_values[string(c)])
	}

	return ip + ":" + decodedPort
}
