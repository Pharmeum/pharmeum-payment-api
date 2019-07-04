package config

import (
	"github.com/caarlos0/env"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	channelID = "pharmeum-channel"
	orgName   = "PharmeumMSP"
	orgAdmin  = "Admin"
)

//Channel read configuration related to pharmeum-channel inside Hyperledger network
type Channel struct {
	ConfigFilePath string `env:"PHARMEUM_CHANNEL_CONFIG_FILE_PATH"`
}

func (c *ConfigImpl) Channel() *channel.Client {
	if c.channelClient != nil {
		return c.channelClient
	}

	c.Lock()
	defer c.Unlock()

	channelConfiguration := &Channel{}
	if err := env.Parse(channelConfiguration); err != nil {
		panic(err)
	}

	configProvider := config.FromFile(channelConfiguration.ConfigFilePath)
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		panic(err)
	}

	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	// Channel client is used to query and execute transactions
	client, err := channel.New(clientChannelContext)
	if err != nil {
		panic(err)
	}

	if client == nil {
		panic("client can't be nil")
	}

	c.channelClient = client

	return c.channelClient
}
