// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
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
	ID    int    `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

type structValidMessage struct {
	regExpr string
	reponse string
}

//MessageReceived :Callback to handle when message received.
func MessageReceived(event Event, opts MessageOpts, msg ReceivedMessage) {
	// log.Println("event:", event, " opt:", opts, " msg:", msg)
	profile, err := mess.GetProfile(opts.Sender.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	validMessages := []structValidMessage{
		{regExpr: "(哈|呵|嘿)", reponse: "笑屁"},
		{regExpr: "(嗨|你好|妳好|您好|哈囉)", reponse: "嗨"},
	}

	var message = fmt.Sprintf(" %s %s : ", profile.FirstName, profile.LastName)
	resp, err := mess.SendSimpleMessage(opts.Sender.ID, message)
	if err != nil {
		fmt.Println(err)
	}
	pages := getPages()
	fmt.Printf("%v", pages[0])
	for _, each := range validMessages {
		valid := regexp.MustCompile(each.regExpr)
		if valid.MatchString(msg.Text) {
			resp, err = mess.SendSimpleMessage(opts.Sender.ID, each.reponse)
			if err != nil {
				fmt.Println(err)
			}
			resp, err = mess.SendSimpleMessage(opts.Sender.ID, pages[0].Title)
			if err != nil {
				fmt.Println(err)
			}
		}
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
