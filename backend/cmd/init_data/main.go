package main

import (
	raw_stocks "backend/internal/features/raw-stocks"
	"log"
)

func main() {

	// Start the data initializer
	initializer, err := raw_stocks.NewDataInitializer()
	if err != nil {
		log.Fatal(err)
	}

	var RawData []raw_stocks.RawItem
	var nextPage string
	var counter int
	for {
		resp, err := initializer.GetData(nextPage)
		if err != nil {
			log.Fatal(err)
		}
		RawData = append(RawData, resp.Items...)
		nextPage = resp.NextPage
		if nextPage == "" || len(resp.Items) == 0 {
			break
		}
		log.Println("Chunk ", counter, "length: ", len(resp.Items), "next page: ", nextPage)
		counter++
	}
	log.Println("Data initialized", RawData)
}
