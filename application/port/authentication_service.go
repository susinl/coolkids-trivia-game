package port

// AuthenticationService is an interface that defines the methods for creating and validating JWT tokens
type AuthenticationService interface {
	// GenerateToken generates a JWT token based on the provided game code and claims
	GenerateToken(gameCode string, claims map[string]interface{}) (string, error)

	// ValidateToken validates a JWT token and returns the game code and any custom claims if the token is valid
	ValidateToken(token string) (string, map[string]interface{}, error)
}
