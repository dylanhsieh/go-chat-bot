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
	"fmt"
	"log"
	"net/http"
	"os"
	"unicode"
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

//MessageReceived :Callback to handle when message received.
func MessageReceived(event Event, opts MessageOpts, msg ReceivedMessage) {
	// log.Println("event:", event, " opt:", opts, " msg:", msg)
	profile, err := mess.GetProfile(opts.Sender.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	//var validID = regexp.MustCompile("^[\u4e00-\u9fa5]+$")
	//var matchResult = validID.MatchString(msg.Text)
	var matchResult = IsChineseChar(msg.Text)
	resp, err := mess.SendSimpleMessage(opts.Sender.ID, fmt.Sprintf("Hello   , %s %s, I don't understand what does %s means %s", profile.FirstName, profile.LastName, msg.Text, matchResult))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", resp)
}

/*
判断字符串是否包含中文字符
*/
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}
