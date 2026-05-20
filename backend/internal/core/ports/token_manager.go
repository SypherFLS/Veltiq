package ports

type TokenManager interface {
	Generate(userID string) (string, error)
	Verify(token string) (string, error)
}