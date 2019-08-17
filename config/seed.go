package config

type Init struct {
    Realms     []string `json:"realms"`
    ReplyRange [2]int   `json:"replyRange"` // only post whose reply count within [min, max) would be downloaded
    PageRange  [2]int   `json:"pageRange"`  // page: [start, end)
    SleepRange [2]int   `json:"sleepRange"` // sleep: [min, max) during download
}
