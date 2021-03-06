package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type Hubstaff struct {
	cookie *http.Cookie
	orgID  string
}

func NewHubstaff(sessionCookie, orgID string) *Hubstaff {
	cookie := &http.Cookie{Name: "_hubstaff_session", Value: sessionCookie}
	return &Hubstaff{
		orgID:  orgID,
		cookie: cookie,
	}
}

const (
	dayFormat  = "Monday"
	timeFormat = "15:04:05"
	weeklyAPI  = "https://app.hubstaff.com/organizations/%s/time_entries.json?filters[view]=weekly"
)

type HubstaffResponse struct {
	TimeEntries []TimeEntry `json:"all_time_entries"`
	TimeZones   []TimeZone  `json:"time_zones"`
}

type TimeEntry struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
	StopTime  time.Time `json:"stop_time"`
	Tracked   int       `json:"tracked"`
	Idle      int       `json:"idle"`
	Billable  bool      `json:"billable"`
}

type TimeZone struct {
	For  string `json:"for"`
	Name string `json:"name"`
}

func (t *Hubstaff) makeAPICall(offset int) (*HubstaffResponse, error) {
	var result *HubstaffResponse

	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(weeklyAPI, t.orgID), nil)

	weekday := time.Now().Weekday()
	if weekday != time.Sunday {
		weekday -= 1
	}
	baseTime := time.Now().AddDate(0, 0, -int(weekday)).AddDate(0, 0, offset*-7)
	endTime := baseTime.AddDate(0, 0, 7)

	q := req.URL.Query()
	q.Set("date", baseTime.Format("2006-01-02"))
	q.Set("date_end", endTime.Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	req.AddCookie(t.cookie)
	fmt.Println(req.URL)

	response, err := client.Do(req)
	if err != nil {
		return result, err
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

func (t *Hubstaff) WeeklyStats(offset int) (string, error) {
	resp, err := t.makeAPICall(offset)
	if err != nil {
		return "", err
	}

	timeLoc, err := getTimeLocation(resp.TimeZones)
	if err != nil {
		return "", err
	}

	return t.renderWeeklyStats(resp.TimeEntries, timeLoc)
}

func (t *Hubstaff) renderWeeklyStats(timeEntries []TimeEntry, timeLoc *time.Location) (string, error) {
	if len(timeEntries) == 0 {
		return "No time entries", nil
	}

	// Total hours
	var total int
	for _, te := range timeEntries {
		total += te.Tracked
	}

	tableString := &strings.Builder{}

	tablePageView := tablewriter.NewWriter(tableString)
	tablePageView.SetBorder(false)
	tablePageView.SetHeader([]string{"Day", "Hours"})

	var (
		dayTotal int
		curDay   time.Time
	)
	curDay = timeEntries[0].StartTime.In(timeLoc)
	for _, te := range timeEntries {
		if te.StartTime.In(timeLoc).Day() == curDay.Day() {
			dayTotal += te.Tracked
			continue
		}

		day := curDay.Format(dayFormat)
		duration := secondsToMinutes(dayTotal)
		tablePageView.Append([]string{day, duration})

		// Reset
		dayTotal = te.Tracked
		curDay = te.StartTime.In(timeLoc)
	}

	// Last day
	day := curDay.Format(dayFormat)
	duration := secondsToMinutes(dayTotal)
	tablePageView.Append([]string{day, duration})

	// Footer
	tablePageView.SetFooter([]string{"", secondsToMinutes(total)})

	tablePageView.Render()

	return tableString.String(), nil
}

func (t *Hubstaff) DailyStats() (string, error) {
	resp, err := t.makeAPICall(0)
	if err != nil {
		return "", err
	}

	timeLoc, err := getTimeLocation(resp.TimeZones)
	if err != nil {
		return "", err
	}

	return t.renderDailyStats(resp.TimeEntries, timeLoc)
}

func (t *Hubstaff) renderDailyStats(timeEntries []TimeEntry, timeLoc *time.Location) (string, error) {
	if len(timeEntries) == 0 {
		return "No time entries", nil
	}

	tableString := &strings.Builder{}

	tablePageView := tablewriter.NewWriter(tableString)
	tablePageView.SetBorder(false)
	tablePageView.SetHeader([]string{"Start", "Stop", "Hours"})

	var dayTotal int
	today := time.Now().In(timeLoc)
	for _, te := range timeEntries {
		if te.StartTime.In(timeLoc).Day() != today.Day() {
			continue
		}

		dayTotal += te.Tracked
		start := te.StartTime.In(timeLoc).Format(timeFormat)
		stop := te.StopTime.In(timeLoc).Format(timeFormat)
		duration := secondsToMinutes(te.Tracked)
		tablePageView.Append([]string{start, stop, duration})
	}

	// Footer
	duration := secondsToMinutes(dayTotal)
	tablePageView.SetFooter([]string{"", "", duration})
	tablePageView.Render()

	return tableString.String(), nil
}

func secondsToMinutes(totalSecs int) string {
	hours := totalSecs / 3600
	minutes := (totalSecs % 3600) / 60
	seconds := totalSecs % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func getTimeLocation(timezones []TimeZone) (*time.Location, error) {
	for _, tz := range timezones {
		if tz.For == "mine" {
			timeLoc, err := time.LoadLocation(tz.Name)
			if err != nil {
				return timeLoc, err
			}
			return timeLoc, err
		}
	}
	return nil, errors.New("timezone not found")
}
