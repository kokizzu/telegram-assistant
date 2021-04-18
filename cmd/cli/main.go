package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/adityathebe/telegram-assistant/services"
)

var (
	SA_API_KEY       string // Simple Analytics API Key
	SA_SITE_NAME     string // Simple Analytics site name
	HUBSTAFF_SESSION string // Hubstaff session cookie
	HUBSTAFF_ORG_ID  string // Hubstaff organization id
)

func main() {
	readKeys()

	services := services.NewService(SA_API_KEY, SA_SITE_NAME, HUBSTAFF_SESSION, HUBSTAFF_ORG_ID)

	sa := flag.Bool("sa", false, "Simple Analytics")
	hsW := flag.Bool("hs", false, "Show hubstaff weekly stats")
	hsD := flag.Bool("hd", false, "Show hubstaff daily stats")
	flag.Parse()

	wg := sync.WaitGroup{}

	if *hsW {
		wg.Add(1)
		go func() {
			stats, err := services.Hubstaff.WeeklyStats()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(stats)
			wg.Done()
		}()
	}

	if *hsD {
		wg.Add(1)
		go func() {
			stats, err := services.Hubstaff.DailyStats()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(stats)
			wg.Done()
		}()
	}

	if *sa {
		wg.Add(1)
		go func() {
			stats, err := services.SimpleAnalytics.DailySummary()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(stats)
			wg.Done()
		}()
	}

	wg.Wait()
}
