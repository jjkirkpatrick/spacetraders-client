package models

type JumpGate struct {
	Symbol     string   `json:"symbol"`
	Connection []string `json:"connections"`
}
