package services

type Services struct {
	Hubstaff        *Hubstaff
	SimpleAnalytics *SimpleAnalytics
	NepseAlpha      *NepseAlpha
}

func NewService(saAPIKey, saSiteName, hubstaffSession, hubstaffOrgID, nepseAlphaCookie, nepseAlphaHolderID string) *Services {
	return &Services{
		SimpleAnalytics: NewSimpleAnalytics(saAPIKey, saSiteName),
		Hubstaff:        NewHubstaff(hubstaffSession, hubstaffOrgID),
		NepseAlpha:      NewNepseAlpha(nepseAlphaCookie, nepseAlphaHolderID),
	}
}
