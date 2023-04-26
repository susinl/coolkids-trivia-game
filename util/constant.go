package util

const (
	ValidateFieldError string = "Invalid Parameters"
)

const (
	ContentType               string = "Content-Type"
	ApplicationJSON           string = "application/json"
	ApplicationUrlEncoded     string = "application/x-www-form-urlencoded"
	AccessControl             string = "Access-Control-Allow-Origin"
	XRequestID                string = "X-Request-ID"
	BasicAuthenticationHeader string = "WWW-Authenticate"
	ReCaptchaTokenHeader      string = "ReCaptcha-Token"
	TokenCtxKey               string = "token"
	ReCaptchaSecretTag        string = "secret"
	ReCaptchaTokenTag         string = "response"
)

const (
	ReadyStatus   string = "ready"
	PendingStatus string = "pending"
	UsedStatus    string = "used"
)

const (
	DateTimeFormat string = "2006-01-02 15:04:05"
)
