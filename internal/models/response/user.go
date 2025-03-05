package response

type LoginChallengeResponse struct {
	Salt      []byte `json:"salt"`
	Challenge []byte `json:"challenge"`
	Version   int    `json:"version"`
}

type UserResponse struct {
	Username         string `json:"username"`
	Email            string `json:"email"`
	Pubkey           []byte `json:"pubkey"`
	EncryptedContent []byte `json:"encryptedContent"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
