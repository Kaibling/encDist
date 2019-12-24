package main

import (
	"net/http"
	"encoding/json"
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/kaibling/encDist/libs"
		"io/ioutil"
	"net/url"
)

func main() {

	//get token
	connectionString := "http://127.0.0.1:8070/token"
	sendUser := &libs.User{Name: "hans",Password:"hanspwd"}
	bytesRepresentation, err := json.Marshal(sendUser)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(connectionString, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	log.Println(string(data))
	token:= string(data)

	//encrypt data
	connectionString = "http://127.0.0.1:8070/encrypt"
	responseData := &libs.PlainDataTransfer{Data: []byte("hans"),Token: token}
	bytesRepresentation, err = json.Marshal(responseData)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err = http.Post(connectionString, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}
    data, err = ioutil.ReadAll(resp.Body)
    identifier := string(data)
    log.Println(identifier)
    

	//decrypt data
	connectionString = "http://127.0.0.1:8070/decrypt"
	postData := url.Values{}
	postData.Add("token", responseData.Token)
	postData.Add("hash", identifier)
	resp, err = http.Post(connectionString, "application/x-www-form-urlencoded; param=value", bytes.NewBufferString(postData.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	data, err = ioutil.ReadAll(resp.Body)
	log.Println(string(data))


}