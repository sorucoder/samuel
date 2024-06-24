package payloads

type CreatePasswordChange struct {
	Email string `json:"email" binding:"required,email"`
}

type FulfillPasswordChange struct {
	NewPassword string `json:"newPassword" binding:"required"`
}
