package engine

import (
	"goproxy/internal/provider"
	"goproxy/pkg/logger"
	"sync"

	"github.com/alitto/pond/v2"
	"github.com/sirupsen/logrus"
)

type FinderArgs struct {
	Type []string
}

type Finder struct {
	Args FinderArgs
}

func NewFinder(args FinderArgs) *Finder {
	return &Finder{
		Args: args,
	}
}
func (f *Finder) Find() []provider.Proxy {
	log := logger.GetLogger()

	log.WithFields(logrus.Fields{
		"type": f.Args.Type,
	}).Debug("Finding proxy servers")

	// Find proxy servers
	providers := []provider.Provider{
		&provider.ProxyListOrg{},
		&provider.XseoIn{},
		&provider.FreeProxyCz{},
		&provider.ProxyListPlus{},
		&provider.ProxyListDownload{},
		&provider.GithubSocks5{},
	}

	var (
		findedProxies = make(map[string]provider.Proxy) // 使用 map 模拟 Set
		mu            sync.Mutex
	)

	pool := pond.NewPool(len(providers))

	var wg sync.WaitGroup
	wg.Add(len(providers))

	for _, p := range providers {
		provider := p
		pool.Submit(func() {
			defer wg.Done()

			proxies, err := provider.FetchProxies(f.Args.Type)
			if err != nil {
				log.WithError(err).Error("Failed to get proxies")
				return
			}

			mu.Lock()
			for _, proxy := range proxies {
				key := proxy.IP + ":" + proxy.Port // 使用 IP 和端口作为唯一键
				if _, exists := findedProxies[key]; !exists {
					log.WithFields(logrus.Fields{
						"provider": p.Name(),
						"ip":       proxy.IP,
						"port":     proxy.Port,
						"type":     proxy.Type,
					}).Debug("Found proxy")

					findedProxies[key] = proxy
				}
			}
			mu.Unlock()
		})
	}

	wg.Wait()
	pool.StopAndWait()

	// 将 map 转为 slice
	result := make([]provider.Proxy, 0, len(findedProxies))
	for _, proxy := range findedProxies {
		result = append(result, proxy)
	}

	log.WithFields(logrus.Fields{
		"count": len(result),
	}).Info("Find proxy servers finished")

	return result
}
