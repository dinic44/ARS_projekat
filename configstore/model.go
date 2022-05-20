package configstore

import (
	"fmt"
	"github.com/google/uuid"
)

type Config struct {
	Id      string            `json:"id"`
	Entries map[string]string `json:"entries"`
	Version string            `json:"version"`
}

const (
	posts = "posts/%s"
	all   = "posts"
)

func generateKey() (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(posts, id), id
}

func constructKey(id string) string {
	return fmt.Sprintf(posts, id)
}
