package main

import "os"

func readKeys() {
	var ok bool

	SA_API_KEY, ok = os.LookupEnv("SA_API_KEY")
	if !ok {
		panic("Simple Analytics API Key is required.")
	}

	SA_SITE_NAME, ok = os.LookupEnv("SA_SITE_NAME")
	if !ok {
		panic("Simple Analytics site name is required")
	}

	HUBSTAFF_SESSION, ok = os.LookupEnv("HUBSTAFF_SESSION")
	if !ok {
		panic("Hubstaff session is required.")
	}

	HUBSTAFF_ORG_ID, ok = os.LookupEnv("HUBSTAFF_ORG_ID")
	if !ok {
		panic("Hubstaff organization id is required.")
	}

	NA_COOKIE, ok = os.LookupEnv("NA_COOKIE")
	if !ok {
		panic("Nepse Alpha Cookie is required.")
	}

	NA_HOLDER_ID, ok = os.LookupEnv("NA_HOLDER_ID")
	if !ok {
		panic("Nepse Alpha Portfolio Hodler ID is required.")
	}
}
