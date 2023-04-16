package winners

type GetWinnersResponse struct {
	Data []*Winner `json:"data,omitempty"`
}

// type GetWinnersData struct {
// 	FullName    string `json:"fullname" example:"John Doe"`
// 	Email       string `json:"email" example:"johndoe@example.com"`
// 	PhoneNumber string `json:"phoneNumber" example:"123-456-7890"`
// 	Code        string `json:"code" example:"ABCD"`
// 	Timestamp   string `json:"timestamp" example:"2023-04-15T13:00:00Z"`
// }
