package engine

import (
	"log"
)

type SimpleEngine struct{}

func (e SimpleEngine) Run(realms ...Request) {
	var requests []Request

	for _, value := range realms {
		requests = append(requests, value)
	}

	itemCount := 0
	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]

		parseResult, err := Worker(r)
		if err != nil {
			continue
		}

		requests = append(requests, parseResult.Requests...)

		for _, item := range parseResult.Items {
			log.Printf("Got item #%d: %v", itemCount, item)
			itemCount++
		}
	}
}
