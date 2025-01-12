package provider

import "testing"

func TestProxyListDownload(t *testing.T) {
	provider := ProxyListDownload{}
	proxies, err := provider.FetchProxies([]string{"https"})
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
