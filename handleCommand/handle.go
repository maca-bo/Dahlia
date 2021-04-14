package command

import (
	"encoding/base64"
	"encoding/json"

	"google.golang.org/api/cloudiot/v1"
)

func (c *Command) Handle(cli *cloudiot.Service, clientID string) error {
	if c.Operate == Unlock || c.Operate == Lock {
		err := c.sendToDevice(cli, clientID)
		if err != nil {
			return err
		}
		return nil
	}
	return &UnimplCommandError{cmd: c.Operate}
}

func (c *Command) sendToDevice(cli *cloudiot.Service, clientID string) error {
	payloadJSON, err := json.Marshal(c)
	if err != nil {
		return err
	}
	req := cloudiot.SendCommandToDeviceRequest{
		BinaryData: base64.StdEncoding.EncodeToString(payloadJSON),
		Subfolder:  "operate",
	}
	resp, err := cli.Projects.Locations.Registries.Devices.SendCommandToDevice(clientID, &req).Do()
	if err != nil {
		return err
	}
	if resp.HTTPStatusCode != 200 {
		return &BadStatusCode{resp.HTTPStatusCode}
	}
	return nil
}
