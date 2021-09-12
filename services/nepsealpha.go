package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
)

const (
	portfolioAPI = "https://nepsealpha.com/tradingusers/load_chartuser_portfolio"
)

type NepseAlpha struct {
	holderID string
	cookie   *http.Cookie
}

func NewNepseAlpha(sessionCookie, holderID string) *NepseAlpha {
	cookie := &http.Cookie{Name: "laravel_session", Value: sessionCookie}
	return &NepseAlpha{
		holderID: holderID,
		cookie:   cookie,
	}
}

type NepseAlphaPortfolioResponse struct {
	Performance string `json:"portfolioPerformance"`
	Tables      string `json:"portfolioTables"`
}

type Portfolio struct {
	Investment   string
	CurrentValue string
	dailyGain    string
	netGain      string
	Items        []PortfolioItem
}

type PortfolioItem struct {
	symbol        string
	quantity      int
	purchasePrice float64
	ltp           float64
	dailyChange   string
	netChange     string
	netProfit     float64
}

func (t *NepseAlpha) makeAPICall(holderID string) (*NepseAlphaPortfolioResponse, error) {
	var result *NepseAlphaPortfolioResponse

	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, portfolioAPI, nil)

	q := req.URL.Query()
	q.Set("holder_id", holderID)
	req.URL.RawQuery = q.Encode()

	req.AddCookie(t.cookie)

	response, err := client.Do(req)
	if err != nil {
		return result, err
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

func (t *NepseAlpha) PortfolioDailySummary() (string, error) {
	p, err := t.Portfolio()
	if err != nil {
		return "", err
	}

	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Total Investment: %s\n", p.Investment))
	s.WriteString(fmt.Sprintf("Current Value: %s\n", p.CurrentValue))
	s.WriteString(fmt.Sprintf("Daily Gain: %s\n", p.dailyGain))
	s.WriteString(fmt.Sprintf("Net Gain: %s\n\n\n", p.netGain))

	// Generate tables
	tableString := &strings.Builder{}
	tablePageView := tablewriter.NewWriter(tableString)
	// tablePageView.SetHeader([]string{"SYM", "QTY", "Daily %", "Net %", "Net"})

	for _, item := range p.Items {
		tablePageView.Append([]string{item.symbol, strconv.Itoa(item.quantity), item.dailyChange, item.netChange, fmt.Sprintf("%.0f", item.netProfit)})
	}

	tablePageView.SetBorder(false)
	tablePageView.SetColumnSeparator("")
	tablePageView.Render()

	s.WriteString(tableString.String())
	return s.String(), nil
}

func (t *NepseAlpha) Portfolio() (Portfolio, error) {
	resp, err := t.makeAPICall(t.holderID)
	if err != nil {
		return Portfolio{}, err
	}

	reader := strings.NewReader(resp.Performance)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return Portfolio{}, err
	}

	var portfolio Portfolio
	portfolio.Investment = strings.TrimSpace(doc.Find("#curentInvestement").Text())
	portfolio.CurrentValue = strings.TrimSpace(doc.Find("#portfolioValue").Text())
	portfolio.dailyGain = strings.TrimSpace(doc.Find("#unrealizedGain").Text())
	portfolio.netGain = strings.TrimSpace(doc.Find("#unrealizedGainAnalyzedResult table tbody tr:nth(1)").Text())

	reader2 := strings.NewReader(resp.Tables)
	doc2, err := goquery.NewDocumentFromReader(reader2)
	if err != nil {
		return Portfolio{}, err
	}

	doc2.Find("#current-portfolio > table > tbody tr").Each(func(i int, s *goquery.Selection) {
		var item PortfolioItem
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			switch i {
			case 1:
				item.symbol = text
			case 3:
				item.dailyChange = text
			case 4:
				item.quantity, _ = strconv.Atoi(text)
			case 5:
				item.purchasePrice, _ = strconv.ParseFloat(text, 64)
			case 6:
				item.ltp, _ = strconv.ParseFloat(text, 64)
			case 9:
				text = strings.ReplaceAll(text, "NPR ", "")
				text = strings.ReplaceAll(text, " )", "")
				data := strings.Split(text, " (")
				if len(data) == 2 {
					item.netChange = data[0]
					item.netProfit, _ = strconv.ParseFloat(data[1], 64)
				}
			}
		})
		portfolio.Items = append(portfolio.Items, item)
	})

	sort.Slice(portfolio.Items, func(a, b int) bool {
		return portfolio.Items[a].netProfit > portfolio.Items[b].netProfit
	})

	return portfolio, nil
}
