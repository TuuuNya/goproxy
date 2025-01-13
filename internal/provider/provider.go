package provider

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

type Proxy struct {
	IP    string
	Port  string
	Type  string
	Delay time.Duration
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

func (p Proxy) Check(max_delay time.Duration) (Proxy, error) {
	// 检查代理是否有效
	if p.IP == "" || p.Port == "" || p.Type == "" {
		return Proxy{}, nil
	}

	startTime := time.Now()

	var checkedProxy Proxy
	var err error

	switch p.Type {
	case "http":
		checkedProxy, err = CheckHttpProxy(p, max_delay)
	case "https":
		checkedProxy, err = CheckHttpsProxy(p, max_delay)
	case "socks5":
		checkedProxy, err = CheckSocks5Proxy(p, max_delay)
	default:
		return Proxy{}, nil
	}

	// 创建一个新的代理对象，并设置延迟字段
	checkedProxy.Delay = time.Since(startTime)

	if err != nil {
		return Proxy{}, err
	}

	return checkedProxy, nil
}

func CheckHttpProxy(proxy Proxy, max_delay time.Duration) (Proxy, error) {
	proxyURL := fmt.Sprintf("http://%s:%s", proxy.IP, proxy.Port)
	parsedProxy, err := url.Parse(proxyURL)
	if err != nil {
		return Proxy{}, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(parsedProxy),
		DialContext: (&net.Dialer{
			Timeout:   max_delay,
			KeepAlive: max_delay,
		}).DialContext,
		TLSHandshakeTimeout: max_delay,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   max_delay,
	}

	resp, err := client.Get("http://httpbin.org/ip")
	if err != nil {
		return Proxy{}, err
	}
	defer resp.Body.Close()

	return proxy, nil
}

func CheckHttpsProxy(proxy Proxy, max_delay time.Duration) (Proxy, error) {
	proxyURL := fmt.Sprintf("https://%s:%s", proxy.IP, proxy.Port)
	parsedProxy, err := url.Parse(proxyURL)
	if err != nil {
		return Proxy{}, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(parsedProxy),
		DialContext: (&net.Dialer{
			Timeout:   max_delay,
			KeepAlive: max_delay,
		}).DialContext,
		TLSHandshakeTimeout: max_delay,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   max_delay,
	}

	resp, err := client.Get("https://httpbin.org/ip")
	if err != nil {
		return Proxy{}, err
	}
	defer resp.Body.Close()

	return proxy, nil
}

func CheckSocks5Proxy(p Proxy, max_delay time.Duration) (Proxy, error) {
	proxyURL := fmt.Sprintf("%s:%s", p.IP, p.Port)

	dialer, err := proxy.SOCKS5("tcp", proxyURL, nil, nil)
	if err != nil {
		return Proxy{}, err
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   max_delay,
	}

	resp, err := client.Get("https://www.cloudflare.com/cdn-cgi/trace")
	if err != nil {
		return Proxy{}, err
	}
	defer resp.Body.Close()

	return p, nil
}

func RemoveProxy(proxies []Proxy, proxy Proxy) []Proxy {
	var newProxies []Proxy
	for _, p := range proxies {
		if p.IP != proxy.IP || p.Port != proxy.Port {
			newProxies = append(newProxies, p)
		}
	}
	return newProxies
}
