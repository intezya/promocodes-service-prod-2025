package user

type CreateUserResponse struct {
	Token string `json:"token"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

type ActivatePromoResponse struct {
	Promo string `json:"promo"`
}
