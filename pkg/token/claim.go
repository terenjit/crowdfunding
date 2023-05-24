package token

type Claim struct {
	Username      string `json:"username"`
	Name          string `json:"name"`
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	RefreshToken  string `json:"refreshToken"`
	Key           string `json:"key"`
	Authorization string `json:"authorization"`
}
