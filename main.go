// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions ad
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var mess = &Messenger{}

func main() {

	port := os.Getenv("PORT")
	log.Println("Server start in port:", port)
	mess.VerifyToken = os.Getenv("TOKEN")
	mess.AccessToken = os.Getenv("TOKEN")
	log.Println("Bot start in token:", mess.VerifyToken)
	mess.MessageReceived = MessageReceived
	http.HandleFunc("/webhook", mess.Handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type Page struct {
	RegExpr  string `json:"regExpr"`
	Response string `json:"response"`
}

//MessageReceived :Callback to handle when message received.
func MessageReceived(event Event, opts MessageOpts, msg ReceivedMessage) {
	// log.Println("event:", event, " opt:", opts, " msg:", msg)
	profile, err := mess.GetProfile(opts.Sender.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	var message = fmt.Sprintf(" %s %s 您好 ", profile.FirstName, profile.LastName)
	resp, err := mess.SendSimpleMessage(opts.Sender.ID, message)
	if err != nil {
		fmt.Println(err)
	}
	pages := getPages()
	matchCount := 0
	for _, each := range pages {
		valid := regexp.MustCompile(each.RegExpr)
		if valid.MatchString(msg.Text) {
			resp, err = mess.SendSimpleMessage(opts.Sender.ID, each.Response)
			matchCount++
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if matchCount <= 0 {
		resp, err = mess.SendSimpleMessage(opts.Sender.ID, "我不懂您在說什麼, 說中文好嗎")
	}

	fmt.Printf("%+v", resp)
}

func getPages() []Page {
	raw, err := ioutil.ReadFile("./messageResponse.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []Page
	json.Unmarshal(raw, &c)
	return c
}

func connectMongo() {
	const (
		Host     = "ds011725.mlab.com:11725"
		Username = "dylan_hsieh"
		Password = "2juxuuux"
		Database = "message"
	)

	session, err := mgo.DialWithInfo(&mgo.DialInfo{})
	if err == nil {
		fmt.Printf("Connected to %v!\n", session.LiveServers())
	}

}
