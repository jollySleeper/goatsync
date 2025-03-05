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

// Check if the user exists in the database & validate the login request
// Generate a token for the user
// Return the token and user information
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

func Login(context *gin.Context) {
	var data request.LoginRequest
	if err := utils.ParseRequest(context, &data); err != nil {
		utils.BadReqError(context, err)
		return
	}

	users := repository.GetUsers()
	user, exists := users[data.Username]
	if !exists {
		utils.UnauthorizedError(context, errors.New("user not found"))
		return
	}

	// Validate login request
	host := context.Request.Host
	err := validateLoginRequest(data, user, "login", host)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token := models.Token{
		ID:        uuid.New(),
		Key:       uuid.New().String(),
		UserID:    user.ID,
		User:      user,
		CreatedAt: time.Now(),
		Expiry:    time.Now().Add(time.Hour * 24 * 7),
	}

	tokens := repository.GetTokens()
	tokens[token.Key] = token

	response := response.LoginResponse{
		Token: token.Key,
		User: response.UserResponse{
			Username:         user.Username,
			Email:            user.Email,
			Pubkey:           user.UserInfo.Pubkey,
			EncryptedContent: user.UserInfo.EncryptedContent,
		},
	}

	if err := utils.SendResponse(context, &response); err != nil {
		utils.InternalError(context, err)
	}
}

func Logout(context *gin.Context) {
	token := context.GetHeader("Authorization")
	if token == "" {
		utils.BadReqError(context, errors.New("missing token from Header"))
		return
	}

	tokens := repository.GetTokens()
	_, exists := tokens[token]
	if !exists {
		utils.UnauthorizedError(context, errors.New("invalid token"))
		return
	}

	delete(tokens, token)
	context.Status(http.StatusNoContent)
}

func ChangePassword(context *gin.Context) {
	var data request.ChangePasswordRequest

	if err := utils.ParseRequest(context, &data); err != nil {
		utils.BadReqError(context, err)
		return
	}

	users := repository.GetUsers()
	user, exists := users[data.Username]
	if !exists {
		utils.UnauthorizedError(context, errors.New("user not found"))
		return
	}

	// Validate old password
	err := bcrypt.CompareHashAndPassword([]byte(user.UserInfo.LoginPubkey), []byte(data.OldPassword))
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect old password"})
		return
	}

	// Validate login request
	host := context.Request.Host
	err = validateLoginRequest(request.LoginRequest{
		Username:  data.Username,
		Challenge: data.Challenge,
		Host:      data.Host,
		Action:    data.Action,
	}, user, "change_password", host)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Hash new password
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash new password"})
		return
	}

	// Update user password
	user.UserInfo.LoginPubkey = newPasswordHash
	users[data.Username] = user

	context.Status(http.StatusNoContent)
}

func SignUp(context *gin.Context) {
	var data request.SignUpRequest

	if err := utils.ParseRequest(context, &data); err != nil {
		utils.BadReqError(context, err)
		return
	}

	user := models.User{
		ID:       uuid.New(),
		Username: data.User.Username,
		Email:    data.User.Email,
		UserInfo: &models.UserInfo{
			ID:               uuid.New(),
			OwnerID:          0,
			Salt:             data.Salt,
			LoginPubkey:      data.Pubkey,
			Pubkey:           data.Pubkey,
			EncryptedContent: data.EncryptedContent,
		},
	}

	users := repository.GetUsers()
	tokens := repository.GetTokens()

	users[user.Username] = user

	token := models.Token{
		ID:        uuid.New(),
		Key:       uuid.New().String(),
		UserID:    user.ID,
		User:      user,
		CreatedAt: time.Now(),
		Expiry:    time.Now().Add(time.Hour * 24 * 7),
	}

	tokens[token.Key] = token

	response := response.LoginResponse{
		Token: token.Key,
		User: response.UserResponse{
			Username:         user.Username,
			Email:            user.Email,
			Pubkey:           user.UserInfo.Pubkey,
			EncryptedContent: user.UserInfo.EncryptedContent,
		},
	}

	if err := utils.SendResponse(context, &response); err != nil {
		utils.InternalError(context, err)
	}
}

// TODO: Compare Implementation of this function
func getEncryptionKey(salt []byte) [32]byte {
	var key [32]byte
	hash := blake2b.Sum256([]byte("your_secret_key"))
	copy(key[:], hash[:])
	return key
}

func validateLoginRequest(data request.LoginRequest, user models.User, expectedAction, host string) error {
	encKey := getEncryptionKey(user.UserInfo.Salt)
	var nonce [24]byte
	length := min(24, len(user.UserInfo.Salt))
	copy(nonce[:length], user.UserInfo.Salt[:length])
	decrypted, ok := secretbox.Open(nil, data.Challenge, &nonce, &encKey)
	if !ok {
		return errors.New("invalid challenge")
	}

	var challengeData map[string]interface{}
	if err := json.Unmarshal(decrypted, &challengeData); err != nil {
		return err
	}

	now := time.Now().Unix()
	if data.Action != expectedAction {
		return errors.New("wrong action")
	} else if now-int64(challengeData["timestamp"].(float64)) > 60 {
		return errors.New("challenge expired")
	} else if challengeData["userId"].(string) != user.ID.String() {
		return errors.New("wrong user")
	} else if !strings.HasPrefix(host, data.Host) {
		return errors.New("wrong host")
	}

	return nil
}

// --------------------- Dashboard URL ---------------------

type DashboardURLResponse struct {
	URL string `json:"url" msgpack:"url"`
}

type CallbackContext struct {
	PathParams map[string]string
	User       *models.User
}

// Add the app settings variable
// Allows to configure the dashboard URL from Config File

// TODO: Implement config file
// func init() {
//     appSettings.DashboardURLFunc = func(context *CallbackContext) string {
//         return fmt.Sprintf("https://dashboard.example.com/users/%s", context.User.Username)
//     }
//     appSettings.UseMsgPack = true // or false for JSON
// }

type AppSettings struct {
	DashboardURLFunc func(context *CallbackContext) string
	UseMsgPack       bool // Configuration flag for response format
}

// Add this with your other global variables
var appSettings = AppSettings{
	DashboardURLFunc: nil,
	UseMsgPack:       true, // Set to false to use JSON
}

// Replace the existing dashboardURL function
func DashboardURL(c *gin.Context) {
	// Get tokenStr from Authorization header
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	// Get username from token
	tokens := repository.GetTokens()
	users := repository.GetUsers()

	token, exists := tokens[tokenStr]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Get user from username
	user, exists := users[token.User.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Check if dashboard URL function is configured
	if appSettings.DashboardURLFunc == nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code":   "not_supported",
			"detail": "This server doesn't have a user dashboard.",
		})
		return
	}

	// Create callback context
	context := &CallbackContext{
		PathParams: make(map[string]string),
		User:       &user,
	}

	// Get URL from dashboard URL function
	url := appSettings.DashboardURLFunc(context)

	// Create response
	response := DashboardURLResponse{
		URL: url,
	}

	if appSettings.UseMsgPack {
		// MessagePack response
		packed, err := codec.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode response"})
			return
		}
		c.Data(http.StatusOK, "application/msgpack", packed)
	} else {
		// JSON response
		c.JSON(http.StatusOK, response)
	}
}
