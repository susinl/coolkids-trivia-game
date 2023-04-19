package winners

type Winners struct {
	Winners []*Winner `json:"winners"`
}

type Winner struct {
	FullName    *string `json:"name" db:"name"`
	PhoneNumber *string `json:"phoneNumber" db:"phone_number"`
	Code        *string `json:"gameCode" db:"game_code"`
	Timestamp   *string `json:"registeredTime" db:"registered_time"`
}
