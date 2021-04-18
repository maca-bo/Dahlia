package Dahlia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	verify "github.com/maca-bo/Dahlia/VerifySlackRequest"
	command "github.com/maca-bo/Dahlia/handleCommand"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudiot/v1"
)

func Slash(w http.ResponseWriter, r *http.Request) {
	setup()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Couldn't read request body: %v\n", err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests accepted", 405)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Couldn't PerseForm", 400)
		log.Fatalf("PerseForm: %v\n", err)
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	result, err := verify.VerifyRequest(r)
	if err != nil {
		log.Fatalf("verifyRequest: %v\n", err)
	}
	if !result {
		log.Fatal("Signature did not match.\n")
	}

	cmd := &command.Command{}
	cmd.UserID = r.Form["user_id"][0]
	ope := strings.TrimPrefix(r.Form["command"][0], "/")
	switch ope {
	case "lock":
		cmd.Operate = command.Lock
	case "unlock":
		cmd.Operate = command.Unlock
	}
	cmd.Argument = nil

	ctx := context.Background()
	httpClient, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	if err != nil {
		log.Fatalf("google.DefaultClient: %v", err)
	}

	cli, err := cloudiot.New(httpClient)
	if err != nil {
		log.Fatalf("cloudiot.New: %v", err)
	}

	clientID := fmt.Sprintf("projects/%v/locations/%v/registries/%v/devices/%v", ProjectID, Region, RegistryID, DeviceID)

	err = cmd.Handle(cli, clientID)
	if err != nil {
		log.Fatalf("cmd.Handle: %v", err)
	}

	msg := &Message{
		ResponseType: ReplyChannel,
		Text:         fmt.Sprintf("<@%s> sent %s command", cmd.UserID, cmd.Operate),
		Attachments:  nil,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(msg); err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}

}
