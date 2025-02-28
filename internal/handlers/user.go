package handlers

import (
	"encoding/json"
	"errors"
	"goatsync/internal/codec"
	"goatsync/internal/models"
	"goatsync/internal/models/request"
	"goatsync/internal/models/response"
	"goatsync/internal/repository"
	"goatsync/pkg/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/secretbox"
)


func LoginChallenge(context *gin.Context) {
	var challenge request.LoginChallengeRequest
	if err := utils.ParseRequest(context, &challenge); err != nil {
		utils.BadReqError(context, err)
		return
	}

	users := repository.GetUsers()

	user, exists := users[challenge.Username]
	if !exists {
		utils.UnauthorizedError(context, errors.New("user not found"))
		return
	}

	// Generate salt and challenge
	salt := user.UserInfo.Salt
	encKey := getEncryptionKey(salt)
	challengeData := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"userId":    user.ID.String(),
	}
	challengeBytes, _ := json.Marshal(challengeData)
	var nonce [24]byte
	length := min(24, len(salt))
	copy(nonce[:length], salt[:length])
	encryptedChallenge := secretbox.Seal(nil, challengeBytes, &nonce, &encKey)

	response := response.LoginChallengeResponse{
		Salt:      salt,
		Challenge: encryptedChallenge,
		Version:   user.UserInfo.Version,
	}

	if err := utils.SendResponse(context, &response); err != nil {
		utils.InternalError(context, err)
	}
}

