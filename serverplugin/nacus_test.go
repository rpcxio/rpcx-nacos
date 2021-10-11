package serverplugin

import (
	"os"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/rpcxio/rpcx-nacos/client"
	"github.com/stretchr/testify/assert"
)

func TestNacos(t *testing.T) {
	defer os.RemoveAll("./log")
	defer os.RemoveAll("./cache")

	clientConfig := constant.ClientConfig{
		TimeoutMs:            10 * 1000,
		ListenInterval:       30 * 1000,
		BeatInterval:         5 * 1000,
		NamespaceId:          "public",
		CacheDir:             "./cache",
		LogDir:               "./log",
		UpdateThreadNum:      20,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: true,
	}

	serverConfig := []constant.ServerConfig{{
		IpAddr: "console.nacos.io",
		Port:   80,
	}}

	r := &NacosRegisterPlugin{
		ServiceAddress: "tcp@127.0.0.1:8972",
		ClientConfig:   clientConfig,
		ServerConfig:   serverConfig,
		Cluster:        "test",
		Group:          "test_group",
	}
	err := r.Start()
	assert.NoError(t, err)
	defer r.Stop()

	err = r.Register("Arith", nil, "")
	assert.NoError(t, err)

	d, err := client.NewNacosDiscovery("Arith", "test", "test_group", clientConfig, serverConfig)
	d.WatchService()
	assert.NoError(t, err)
	pairs := d.GetServices()
	assert.Equal(t, 1, len(pairs))
	d.Close()
}
