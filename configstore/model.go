package configstore

import (
	"fmt"
	"github.com/google/uuid"
)

type SingleConfig struct {
	Id      string            `json:"id"`
	Entries map[string]string `json:"entries"`
	Version string            `json:"version"`
}

type GroupConfig struct {
	Id          string              `json:"id"`
	GroupConfig []map[string]string `json:"configs"`
	Version     string              `json:"version"`
}

const (
	singleConfigAll     = "singleConfigs"
	singleConfigId      = "singleConfig/%s"
	singleConfigVersion = "singleConfig/%s/%s"

	groupConfigAll     = "groupConfigs"
	groupConfigId      = "group/%s"
	groupConfigVersion = "group/%s/%s"
)

func generateSingleConfigKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(singleConfigVersion, id, ver), id
}

func constructSingleConfigKey(id string, ver string) string {
	return fmt.Sprintf(singleConfigVersion, id, ver)
}

func constructSingleConfigIdKey(id string) string {
	return fmt.Sprintf(singleConfigId, id)
}

func generateGroupConfigKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(groupConfigVersion, id, ver), id
}

func constructGroupConfigKey(id string, ver string) string {
	return fmt.Sprintf(groupConfigVersion, id, ver)
}

func constructGroupConfigIdKey(id string) string {
	return fmt.Sprintf(groupConfigId, id)
}
