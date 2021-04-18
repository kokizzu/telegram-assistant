package services

type Services struct {
	Hubstaff        *Hubstaff
	SimpleAnalytics *SimpleAnalytics
}

func NewService(saAPIKey, saSiteName, hubstaffSession, hubstaffOrgID string) *Services {
	return &Services{
		SimpleAnalytics: NewSimpleAnalytics(saAPIKey, saSiteName),
		Hubstaff:        NewHubstaff(hubstaffSession, hubstaffOrgID),
	}
}
