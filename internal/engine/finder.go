package engine

import (
	"goproxy/internal/provider"
	"goproxy/pkg/logger"

	"github.com/sirupsen/logrus"
)

type FinderArgs struct {
	MaxDelay int
	Type     []string
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
		"max_delay": f.Args.MaxDelay,
		"type":      f.Args.Type,
	}).Info("Finding proxy servers")

	// Find proxy servers
	providers := []provider.Provider{
		&provider.ProxyListOrg{},
		&provider.XseoIn{},
		&provider.FreeProxyCz{},
		&provider.ProxyListPlus{},
	}
	finded_proxies := []provider.Proxy{}

	for _, p := range providers {
		proxies, err := p.FetchProxies(f.Args.Type)
		if err != nil {
			log.WithError(err).Error("Failed to get proxies")
			continue
		}

		for _, proxy := range proxies {
			log.WithFields(logrus.Fields{
				"provider": p.Name(),
				"ip":       proxy.IP,
				"port":     proxy.Port,
				"type":     proxy.Type,
			}).Info("Found proxy")

			finded_proxies = append(finded_proxies, proxy)
		}
	}

	log.WithFields(logrus.Fields{
		"count": len(finded_proxies),
	}).Info("Find proxy servers finished")

	return finded_proxies
}
