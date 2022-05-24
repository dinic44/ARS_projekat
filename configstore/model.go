package configstore

import (
	"fmt"
	"github.com/google/uuid"
)

type SingleConfig struct {
	Id      string            `json:"id"`
	Version string            `json:"version"`
	Entries map[string]string `json:"entries"`
}

type GroupConfig struct {
	Id          string              `json:"id"`
	GroupConfig []map[string]string `json:"configs"`
	Version     string              `json:"version"`
}

const (
	singleConfigId  = "singleConfig/%s"
	singleConfig    = "singleConfig/%s/%s"
	singleConfigAll = "singleConfigs"

	/*groupConfigAll     = "groupConfigs"
	groupConfigId      = "groupConfig/%s"
	groupConfig 	   = "groupConfig/%s/%s"*/
)

func generateSingleConfigKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(singleConfig, id, ver), id
}

func constructSingleConfigKey(id string, ver string) string {
	return fmt.Sprintf(singleConfig, id, ver)
}

func constructSingleConfigIdKey(id string) string {
	return fmt.Sprintf(singleConfigId, id)
}

/*func generateGroupConfigKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(groupConfigVersion, id, ver), id
}

func constructGroupConfigKey(id string, ver string) string {
	return fmt.Sprintf(groupConfigVersion, id, ver)
}

func constructGroupConfigIdKey(id string) string {
	return fmt.Sprintf(groupConfigId, id)
}*/
