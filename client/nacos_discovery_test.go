package client

import (
	"os"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/stretchr/testify/assert"
)

func TestNacosDiscovery(t *testing.T) {
	defer os.RemoveAll("./log")
	defer os.RemoveAll("./cache")

	clientConfig := constant.ClientConfig{
		NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", // namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		CacheDir:            "./cache",
		LogDir:              "./log",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	serverConfig := []constant.ServerConfig{{
		IpAddr: "console.nacos.io",
		Port:   80,
	}}

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfig,
		},
	)
	assert.NoError(t, err)

	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "Arith",
		ClusterName: "test",
		GroupName:   "test_group",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})
	assert.NoError(t, err)
	assert.True(t, success)

	time.Sleep(100 * time.Millisecond)

	d, err := NewNacosDiscovery("Arith", "test", "test_group", clientConfig, serverConfig)
	d.WatchService()
	assert.NoError(t, err)
	pairs := d.GetServices()
	assert.Equal(t, 1, len(pairs))
	d.Close()
}
