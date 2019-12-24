package main

import (
	"github.com/kaibling/encDist/publisher"
	log "github.com/sirupsen/logrus"
	"github.com/kaibling/encDist/libs"
)

func main() {

	log.SetLevel(log.DebugLevel)
	cliArguments := libs.ParseArguments()
	config := libs.ParseConfigurationFile(cliArguments["configFilePath"])
	publisher := publisher.NewPublisher(config)
	publisher.StartServer()
}