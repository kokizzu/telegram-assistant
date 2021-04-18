package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type SimpleAnalytics struct {
	apiKey string
	apiURL string
}

func NewSimpleAnalytics(apiKey, siteName string) *SimpleAnalytics {
	return &SimpleAnalytics{
		apiKey: apiKey,
		apiURL: fmt.Sprintf("https://simpleanalytics.com/%s.json", siteName),
	}
}

// Auto generated with https://mholt.github.io/json-to-go/
type APIResponse struct {
	Ok        bool      `json:"ok"`
	Docs      string    `json:"docs"`
	Info      string    `json:"info"`
	Hostname  string    `json:"hostname"`
	URL       string    `json:"url"`
	Path      string    `json:"path"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Version   int       `json:"version"`
	Timezone  string    `json:"timezone"`
	Histogram []struct {
		Date      string `json:"date"`
		Pageviews int    `json:"pageviews"`
		Visitors  int    `json:"visitors"`
	} `json:"histogram"`
	Pages []struct {
		Value     string `json:"value"`
		Pageviews int    `json:"pageviews"`
		Visitors  int    `json:"visitors"`
	} `json:"pages"`
	Countries []struct {
		Value     string `json:"value"`
		Pageviews int    `json:"pageviews"`
		Visitors  int    `json:"visitors"`
	} `json:"countries"`
	GeneratedInMs int `json:"generated_in_ms"`
}

func (t *SimpleAnalytics) callAPI(options map[string]string) (*APIResponse, error) {
	r := &APIResponse{}

	req, err := http.NewRequest("GET", t.apiURL, nil)
	if err != nil {
		return r, err
	}

	q := req.URL.Query()
	q.Add("version", "5")
	q.Add("info", "false")
	q.Add("fields", "histogram,pages,countries")
	q.Add("timezone", "Asia/Kathmandu")
	for k, v := range options {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Api-Key", t.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return r, err
	}

	return r, nil
}

func (t *SimpleAnalytics) DailySummary() (string, error) {
	options := make(map[string]string)
	options["start"] = getTime()

	apiResp, err := t.callAPI(options)
	if err != nil {
		return "", err
	}

	return t.renderDailySummary(apiResp)
}

func (t *SimpleAnalytics) renderDailySummary(r *APIResponse) (string, error) {
	tableString := &strings.Builder{}
	tableString.WriteString("<pre>")

	// Views by Page
	tablePageView := tablewriter.NewWriter(tableString)
	tablePageView.SetHeader([]string{"Page", "Views", "Visitors"})
	tablePageView.SetBorder(false)

	for _, p := range r.Pages {
		tablePageView.Append([]string{p.Value, fmt.Sprint(p.Pageviews), fmt.Sprint(p.Visitors)})
	}

	for _, h := range r.Histogram {
		tablePageView.SetFooter([]string{"Total", fmt.Sprint(h.Pageviews), fmt.Sprint(h.Visitors)})
	}

	tablePageView.Render()
	tableString.WriteString("</pre>\n\n<pre>")

	// Views by Country
	tableCountries := tablewriter.NewWriter(tableString)
	tableCountries.SetHeader([]string{"Country", "Views", "Visitors"})
	tableCountries.SetBorder(false)

	for _, c := range r.Countries {
		tableCountries.Append([]string{c.Value, fmt.Sprint(c.Pageviews), fmt.Sprint(c.Visitors)})
	}
	tableCountries.Render()

	tableString.WriteString("</pre>")

	return tableString.String(), nil
}

// getTime returns date in format : 2021-04-04
func getTime() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}
