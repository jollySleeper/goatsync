package request

type LoginChallengeRequest struct {
	Username string `json:"username" binding:"required"`
}

type LoginRequest struct {
	Username  string `json:"username"`
	Challenge []byte `json:"challenge"` // base64 encoded
	Host      string `json:"host"`
	Action    string `json:"action"`
}

type ChangePasswordRequest struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Challenge   []byte `json:"challenge"` // base64 encoded
	Host        string `json:"host"`
	Action      string `json:"action"`
}

type SignUpRequest struct {
	User             UserSignUpRequest `json:"user"`
	Salt             []byte            `json:"salt"`
	LoginPubkey      []byte            `json:"loginPubkey"`
	Pubkey           []byte            `json:"pubkey"`
	EncryptedContent []byte            `json:"encryptedContent"`
}

type UserSignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
