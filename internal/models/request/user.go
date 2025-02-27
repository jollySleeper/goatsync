package request

type LoginChallengeRequest struct {
	Username string `json:"username" binding:"required"`
}
