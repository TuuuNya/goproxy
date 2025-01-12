package provider

import "testing"

func TestProxyListPlus(t *testing.T) {
	provider := ProxyListPlus{}
	proxies, err := provider.FetchProxies([]string{"http", "https", "socks5"})
	if err != nil {
		t.Errorf("Error fetching proxies: %v", err)
	}
	if len(proxies) == 0 {
		t.Errorf("No proxies fetched")
	}

	for _, proxy := range proxies {
		t.Logf("Proxy: %s:%s (%s)", proxy.IP, proxy.Port, proxy.Type)
	}
}
