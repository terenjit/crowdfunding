package token

type Claim struct {
	Username     string `json:"username"`
	UserID       string `json:"userId"`
	RefreshToken string `json:"refreshToken"`
	Key          string `json:"key"`
}
