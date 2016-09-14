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
	"regexp"
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
	var validLaugh = regexp.MustCompile("(哈|呵|嘿)")
	//var validGreeting = regexp.MustCompile("(嗨|你好|妳好|您好)")
	var validMessage []string
	var messageResponse []string
	validMessage[0] = "(哈|呵|嘿)"
	validMessage[1] = "(嗨|你好|妳好|您好)"
	messageResponse[0] = "笑屁"
	messageResponse[1] = "嗨"
	var matchLaughResult = validLaugh.MatchString(msg.Text)
	//var matchResult = IsChineseChar(msg.Text)
	//var message = fmt.Sprintf("Hello   , %s %s", profile.FirstName, profile.LastName)
	//if matchResult {
	//message = fmt.Sprintf("嗨   , %s %s", profile.FirstName, profile.LastName)
	//}
	var message = fmt.Sprintf(" %s %s : ", profile.FirstName, profile.LastName)
	resp, err := mess.SendSimpleMessage(opts.Sender.ID, message)
	if err != nil {
		fmt.Println(err)
	}
	for key, value := range validMessage {
		resp, err = mess.SendSimpleMessage(opts.Sender.ID, key)
		resp, err = mess.SendSimpleMessage(opts.Sender.ID, value)
	}
	if matchLaughResult {
		resp, err = mess.SendSimpleMessage(opts.Sender.ID, "笑屁")
		if err != nil {
			fmt.Println(err)
		}
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
