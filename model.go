package main

type Config struct {
	Entries map[string]string `json:"entries"`
	Version string            `json:"version"`
}
