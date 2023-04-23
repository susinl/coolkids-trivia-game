package winners

type GetWinnersResponse struct {
	Data []*Winner `json:"data,omitempty"`
}

// type GetWinnersData struct {
// 	Name    string `json:"fullname" example:"John Doe"`
// 	PhoneNumber string `json:"phoneNumber" example:"123-456-7890"`
// 	Code        string `json:"code" example:"ABCD"`
// 	Timestamp   string `json:"timestamp" example:"2023-04-15T13:00:00Z"`
// }

type GetQuotaResponse struct {
	Quota int `json:"quota" example:"0"`
}

type UpdateQuotaRequest struct {
	NewQuota int `json:"newQuota" example:"0"`
}
