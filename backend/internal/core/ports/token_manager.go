package ports

type TokenPair struct {
	AccessToken string
	RefreshToken string
}

type TokenManager interface {
	GeneratePair(userID string) (TokenPair, error)
	VerifyAccess(token string) (userID string, err error)
	VerifyRefresh(token string) (userID string, err error)
}
