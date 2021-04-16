package Dahlia

import "os"

var (
	ProjectID  string
	Region     string
	RegistryID string
	DeviceID   string
)

func setup() {
	ProjectID = os.Getenv("PROJECT_ID")
	Region = os.Getenv("REGION")
	RegistryID = os.Getenv("REGISTRY_ID")
	DeviceID = os.Getenv("DEVICE_ID")
}
