package schemas

type RegistrationSchema struct {
	Login    string `json:"login" binding:"required,max=200"`
	Password string `json:"password" binding:"required"`
}

type RefreshAccessTokenSchema struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}