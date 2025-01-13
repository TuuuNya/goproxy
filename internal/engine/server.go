package engine

import (
	"context"
	"fmt"
	"goproxy/internal/provider"
	"goproxy/pkg/logger"
	"math"
	"net"
	"time"

	"github.com/armon/go-socks5"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

type ServerArgs struct {
	Type          string
	Port          int
	MaxDelay      time.Duration
	CheckPoolSize int
}

type Server struct {
	Args      ServerArgs
	ProxyPool []provider.Proxy
}

func NewServer(args ServerArgs) *Server {
	return &Server{
		Args: args,
	}
}

func (s *Server) Serve() error {
	log := logger.GetLogger()

	go s.updateProxyPool()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Args.Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	log.WithFields(logrus.Fields{
		"port": s.Args.Port,
	}).Infof("Proxy server listening on port %d", s.Args.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.WithError(err).Error("Failed to accept connection")
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(clientConn net.Conn) {
	log := logger.GetLogger()
	defer clientConn.Close()

	// 随机选择一个代理
	proxyInstance := s.getProxyWithMinDelay()
	if proxyInstance.IP == "" {
		log.Warn("No valid proxy selected")
		return
	}

	proxyAddr := fmt.Sprintf("%s:%s", proxyInstance.IP, proxyInstance.Port)

	// 创建 SOCKS5 客户端 Dialer
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		log.WithError(err).Error("Failed to create socks5 proxy dialer")
		return
	}

	wrappedDial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}

	// 创建 SOCKS5 服务器配置
	conf := &socks5.Config{
		Dial: wrappedDial,
	}

	// 创建 SOCKS5 服务器实例
	server, err := socks5.New(conf)
	if err != nil {
		log.WithError(err).Error("Failed to create socks5 server")
		return
	}

	// 处理客户端连接
	if err := server.ServeConn(clientConn); err != nil {
		log.WithError(err).Error("Failed to serve socks5 connection")
	}
}

func (s *Server) getProxyWithMinDelay() provider.Proxy {
	log := logger.GetLogger()

	// 确保 ProxyPool 不为空
	if len(s.ProxyPool) == 0 {
		log.Warn("No valid proxies")
		return provider.Proxy{}
	}

	// 寻找延迟最低的 Proxy
	var minLatencyProxy provider.Proxy
	minLatency := time.Duration(math.MaxInt64) // 初始化为最大值

	for _, proxy := range s.ProxyPool {
		if proxy.Delay < minLatency {
			minLatency = proxy.Delay
			minLatencyProxy = proxy
		}
	}

	log.WithFields(logrus.Fields{
		"ip":    minLatencyProxy.IP,
		"port":  minLatencyProxy.Port,
		"delay": minLatencyProxy.Delay,
	}).Debug("Selected min delay proxy")

	return minLatencyProxy
}

func (s *Server) updateProxyPool() {
	log := logger.GetLogger()
	for {
		validProxies := s.getValidProxy()

		s.ProxyPool = validProxies

		log.WithFields(logrus.Fields{
			"count": len(s.ProxyPool),
		}).Info("Got valid proxies")

		time.Sleep(1 * time.Minute)

	}
}

func (s *Server) getValidProxy() []provider.Proxy {
	proxies := s.getProxy()
	Checker := NewChecker(CheckerArgs{
		SourceProxies: proxies,
		MaxDelay:      s.Args.MaxDelay,
		PoolSize:      s.Args.CheckPoolSize,
	})
	return Checker.Check()
}

func (s *Server) getProxy() []provider.Proxy {
	log := logger.GetLogger()
	log.Info("Getting proxies...")
	finder := NewFinder(FinderArgs{Type: []string{s.Args.Type}})
	return finder.Find()
}
