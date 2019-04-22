package config

type Init struct {
	Seeds        []string `json:"seeds"`
	TopicPerPage int      `json:"topicPerPage"`
	StartPageNum int      `json:"startPageNum"`
	EndPageNum   int      `json:"endPageNum"`
}
