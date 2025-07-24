package dto

type LoginDto struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
