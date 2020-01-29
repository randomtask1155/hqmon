package main 

import (
	"time"
	"os"
	"net/http"
	"log"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/randomtask1155/hqserver/device"
	"crypto/tls"
)

var (
	logger         *log.Logger
	monEndpoint string 
	sendGridAPIKey string
	receiverName string 
	receiverEmail string
	senderName string
	senderEmail string 
)

func getStatus() ([]device.DeviceStatus,error) {
	mon := make([]device.DeviceStatus,0)

	tr := &http.Transport{
		MaxIdleConns:       -1,
		IdleConnTimeout:    30 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr,
			Timeout: 30 * time.Second,
	}

	r, err := client.Get(monEndpoint)
	if err != nil {
		logger.Println(err)
		return mon, err 
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Println(err)
		return mon, err
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(bytes.NewBuffer(body))
	err = decoder.Decode(&mon)
	if err != nil {
		logger.Println(err)
		return mon, err
	}
	return mon, nil
}

func reportStatus(mon []device.DeviceStatus, monErr error) {
	messBody := ""
	
	if monErr != nil {
		messBody = fmt.Sprintf("unable to get status from hq: %s", monErr)
		err := SendGmail(senderName, senderEmail, receiverName, receiverEmail, "Status Report", messBody, messBody)
		if err != nil {
			logger.Println(err)
		}
		return 
	}
	
	sendMail := false
	for i := range mon {
		if mon[i].Status == device.UnHealthyStatus {
			sendMail = true
		}
		messBody += fmt.Sprintf(`<p style="border:3px; border-style:solid; border-color:#5bc0de; padding: 1em;"><b>%s = %s</b></p>`, mon[i].Name, mon[i].Status)
		for e := range mon[i].Events {
			messBody += fmt.Sprintf("<p>%s</p>",mon[i].Events[e])
		}
	}

	if sendMail {
		logger.Println("sending report")
		err := SendGmail(senderName, senderEmail, receiverName, receiverEmail, "Status Report", messBody, messBody)
		if err != nil {
			logger.Println(err)
		}
	}
}

func main() {
	logger = log.New(os.Stdout, "logger: ", log.Ldate|log.Ltime|log.Lshortfile)
	monEndpoint = os.Getenv("MON_ENDPOINT")
	sendGridAPIKey = os.Getenv("SENDGRID_KEY")
	receiverName = os.Getenv("RECEIVER_NAME")
	receiverEmail = os.Getenv("RECEIVER_EMAIL")
	senderName = os.Getenv("SENDER_NAME")
	senderEmail = os.Getenv("SENDER_EMAIL")
	
	for {
		mon, err := getStatus()
		if err != nil {
			logger.Println(err)
		}
		reportStatus(mon, err)
		time.Sleep(2 * time.Hour)
	}

}