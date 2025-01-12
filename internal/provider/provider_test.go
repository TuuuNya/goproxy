package provider

import (
	"testing"
	"time"
)

func TestCheckHttpProxy(t *testing.T) {
	proxy := Proxy{
		IP:   "127.0.0.1",
		Port: "7890",
		Type: "http",
	}
	checked_proxy, err := proxy.Check(10 * time.Second)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	t.Logf("Checked proxy: %v", checked_proxy)

}

func TestCheckHttpsProxy(t *testing.T) {
	proxy := Proxy{
		IP:   "114.38.132.17",
		Port: "8888",
		Type: "https",
	}
	checked_proxy, err := proxy.Check(10 * time.Second)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	t.Logf("Checked proxy: %v", checked_proxy)

}

func TestCheckSocks5Proxy(t *testing.T) {
	proxy := Proxy{
		IP:   "132.148.167.243",
		Port: "44970",
		Type: "socks5",
	}
	checked_proxy, err := proxy.Check(10 * time.Second)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	t.Logf("Checked proxy: %v", checked_proxy)

}
