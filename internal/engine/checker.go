package engine

import (
	"fmt"
	"goproxy/internal/provider"
	"time"

	"github.com/panjf2000/ants/v2"
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
	// using ants pool to check proxies concurrently
	pool, _ := ants.NewPool(c.Args.PoolSize)
	defer pool.Release()

	// check proxies
	results := make(chan provider.Proxy, len(c.Args.SourceProxies))
	for _, proxy := range c.Args.SourceProxies {
		proxy := proxy
		_ = pool.Submit(func() {
			checked_proxy, err := proxy.Check(c.Args.MaxDelay)
			if err != nil {
				return
			}

			results <- checked_proxy
		})
	}

	for i := 0; i < len(c.Args.SourceProxies); i++ {
		proxy_results := <-results
		checked_proxies = append(checked_proxies, proxy_results)
		fmt.Println("Alive proxy: ", proxy_results)
	}

	return checked_proxies
}
