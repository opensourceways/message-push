package kafka

type ConsumeConfig struct {
	Topic string `json:"topic"  required:"true"`
	Group string `json:"group"  required:"true"`
}
