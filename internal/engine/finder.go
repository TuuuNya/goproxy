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
		finded_proxies []provider.Proxy
		mu             sync.Mutex
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
				log.WithFields(logrus.Fields{
					"provider": p.Name(),
					"ip":       proxy.IP,
					"port":     proxy.Port,
					"type":     proxy.Type,
				}).Debug("Found proxy")

				finded_proxies = append(finded_proxies, proxy)
			}
			mu.Unlock()
		})
	}

	wg.Wait()
	pool.StopAndWait()

	log.WithFields(logrus.Fields{
		"count": len(finded_proxies),
	}).Info("Find proxy servers finished")

	return finded_proxies
}
