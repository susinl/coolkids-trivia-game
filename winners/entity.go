package winners

type Winners struct {
	Winners []*Winner `json:"winners"`
}

type Winner struct {
	Name        *string `json:"name" db:"name"`
	PhoneNumber *string `json:"phoneNumber" db:"phone_number"`
	Code        *string `json:"code" db:"code"`
	Timestamp   *string `json:"registeredTime" db:"registered_time"`
	RecentTime  *string `json:"recentTime" db:"recent_time"`
}
