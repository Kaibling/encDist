package publisher

import (
	log "github.com/sirupsen/logrus"
	"github.com/kaibling/encDist/libs"
	//"github.com/rs/xid"
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	//"io/ioutil"
)

type Publisher struct {
	configuration *libs.Configuration
}

func NewPublisher(configuration *libs.Configuration ) *Publisher{

	returnPublisher := new(Publisher)
	returnPublisher.configuration = configuration
	log.Debug("Loaded Config")
	log.Debug(returnPublisher.configuration)
	libs.SQLiteInitPshDB(configuration.DBpath)
	return returnPublisher
}


func (Publisher *Publisher) StartServer() {

	libs.SQLiteInitPshDB(Publisher.configuration.DBpath)
    http.HandleFunc("/publish", Publisher.publishHandler)
    http.HandleFunc("/data", Publisher.dataHandler)

	log.Info("server started on Port " + Publisher.configuration.BindingPort )
	http.ListenAndServe(":"+Publisher.configuration.BindingPort, nil)

}

func (Publisher *Publisher) publishHandler(w http.ResponseWriter, r *http.Request) {

	//get data
	var responseData libs.CryptoDataTransfer
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
    json.Unmarshal([]byte(buf.String()), &responseData)

	//save data
	guid := libs.SQLiteaddPublishData(Publisher.configuration.DBpath,responseData.CryptoData)
	fmt.Fprintf(w,guid)
	libs.SQLiteGetALLPublishData(Publisher.configuration.DBpath)
}

func (Publisher *Publisher) dataHandler(w http.ResponseWriter, r *http.Request) {

    err := r.ParseForm()
	if err != nil {
		panic(err)
	}
    hash := r.Form.Get("hash")
    requestedData := libs.SQLiteGetPublishedData(Publisher.configuration.DBpath,hash)
	fmt.Fprintf(w,string(requestedData))

}