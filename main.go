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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
)

var mess = &Messenger{}

func main() {

	const (
		Host       = "ds011725.mlab.com:11725"
		Username   = "dylan_hsieh"
		Password   = "2juxuuux"
		Database   = "message"
		Collection = "messageResponse"
	)
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{Host},
		Username: Username,
		Password: Password,
		Database: Database,
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		},
	})
	if err != nil {
		panic(err)
	}
	if err == nil {
		fmt.Printf("哈哈哈 Connected to %v!\n", session.LiveServers())
	}
	//defer session.Close()

	//coll := session.DB(Database).C(Collection)
	//player := "超可愛"
	//gamesWon, err := coll.Find(bson.M{"response": player}).Count()
	//if err != nil {
	//panic(err)
	//}

	//fmt.Printf("%s has won %d games.\n", player, gamesWon)
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
		resp, err = mess.SendSimpleMessage(opts.Sender.ID, "我不懂您在說什麼")
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

}
