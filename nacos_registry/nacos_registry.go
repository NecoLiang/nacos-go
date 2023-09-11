// nacos_registry/nacos_registry.go

package nacos_registry

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"nacos-go/config"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type NacosRegistry struct {
	client  naming_client.INamingClient
	service vo.RegisterInstanceParam
}

func NewNacosRegistry() (*NacosRegistry, error) {
	appConfig, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NamespaceId:         appConfig.Go.Nacos.Namespace,
		AppName:             appConfig.Go.Service.Name,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
		//Endpoint:    appConfig.Go.Nacos.ServerIP,
		// 填充其他字段...
	}

	serverConfigs := []constant.ServerConfig{
		{IpAddr: appConfig.Go.Nacos.ServerIP, Port: 8848},
		// 填充其他服务器配置...
	}

	nacosClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return nil, err
	}

	// 获取本地 IP 地址
	localIP, err := GetLocalIP()
	if err != nil {
		return nil, err
	}
	service := vo.RegisterInstanceParam{
		Ip:          localIP,
		Port:        uint64(appConfig.Go.Service.Port),
		ServiceName: appConfig.Go.Service.Name,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		// 填充其他字段...
	}

	return &NacosRegistry{
		client:  nacosClient, // 使用实例化后的 NamingClient
		service: service,
	}, nil
}

func (nr *NacosRegistry) RegisterService() error {

	_, err := nr.client.RegisterInstance(nr.service)
	if err != nil {
		return err
	}
	return nil
}

func (nr *NacosRegistry) DeregisterService() error {
	deregisterParam := vo.DeregisterInstanceParam{
		Ip:          nr.service.Ip,
		Port:        nr.service.Port,
		ServiceName: nr.service.ServiceName,
		Ephemeral:   true,
	}

	_, err := nr.client.DeregisterInstance(deregisterParam)
	if err != nil {
		return err
	}
	return nil
}

func (nr *NacosRegistry) Run() {
	// 监听中断信号，用于优雅关闭服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := nr.DeregisterService(); err != nil {
		log.Println("Failed to deregister service:", err)
	}
	log.Println("Service deregistered successfully.")

}

// GetLocalIP 获取本地 IP 地址
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	// 找到非回环地址的第一个 IPv4 或 IPv6 地址
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil // 返回 IPv4 地址
			}
			return ipnet.IP.String(), nil // 返回 IPv6 地址
		}
	}

	return "", errors.New("未找到本地 IP 地址")
}
