package engine

import (
	"goproxy/internal/provider"
	"time"

	"github.com/alitto/pond/v2"
)

type CheckerArgs struct {
	MaxDelay      time.Duration
	SourceProxies []provider.Proxy
	PoolSize      int
}

type Checker struct {
	Args CheckerArgs
}

func NewChecker(args CheckerArgs) *Checker {
	return &Checker{
		Args: args,
	}
}

func (c *Checker) Check() (checked_proxies []provider.Proxy) {
	pool := pond.NewPool(c.Args.PoolSize)

	for _, proxy := range c.Args.SourceProxies {
		pool.Submit(func() {
			valid_proxy, err := proxy.Check(c.Args.MaxDelay)
			if err != nil {
				return
			}
			checked_proxies = append(checked_proxies, valid_proxy)
		})
	}

	pool.StopAndWait()
	return checked_proxies
}
