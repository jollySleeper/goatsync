package response

type LoginChallengeResponse struct {
	Salt      []byte `json:"salt"`
	Challenge []byte `json:"challenge"`
	Version   int    `json:"version"`
}
