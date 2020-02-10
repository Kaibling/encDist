package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/kaibling/encDist/tokenizer"
	"github.com/kaibling/encDist/libs"
	)


func main() {
	log.SetLevel(log.DebugLevel)
	//data := []byte("The s geht wordsuper see thingy")
	cliArguments := libs.ParseArguments()
	config := libs.ParseConfigurationFile(cliArguments["configFilePath"])
	tokenizer := tokenizer.NewTokenizer(config)
	tokenizer.StartServer()
}