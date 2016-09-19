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
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"regexp"
)

var mess = &Messenger{}

type MessageValidResponse struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	RegExpr  string        `bson:"regExpr"`
	Response string        `bson:"response"`
}

const (
	Host       = "ds011725.mlab.com:11725"
	Username   = "dylan_hsieh"
	Password   = "2juxuuux"
	Database   = "message"
	Collection = "messageResponse"
)

var results []MessageValidResponse
var collectionConnection = nil

func main() {

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{Host},
		Username: Username,
		Password: Password,
		Database: Database,
	})
	if err == nil {
		fmt.Printf("Connected to %v!\n", session.LiveServers())
		defer session.Close()
		coll := session.DB(Database).C(Collection)
		collectionConnection = coll
		fmt.Println("Mongo Message: ", collectionConnection)
		var result []MessageValidResponse
		err := coll.Find(bson.M{}).All(&result)
		results = result

		if err == nil {
			fmt.Println("Mongo Message: ", result)
			port := os.Getenv("PORT")
			log.Println("Server start in port:", port)
			mess.VerifyToken = os.Getenv("TOKEN")
			mess.AccessToken = os.Getenv("TOKEN")
			log.Println("Bot start in token:", mess.VerifyToken)
			mess.MessageReceived = MessageReceived
			http.HandleFunc("/webhook", mess.Handler)
			log.Fatal(http.ListenAndServe(":"+port, nil))
		} else {
			log.Println("read fail", err)
		}
	}
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
	matchCount := 0
	for _, each := range results {
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

func connectMongo() {

}
