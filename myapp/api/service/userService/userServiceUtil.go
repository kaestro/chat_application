// myapp/api/service/userService/userServiceUtil.go
package userService

import (
	"errors"
	"fmt"
	"myapp/api/models"
	"myapp/internal/db"
	"myapp/internal/password"
	"myapp/internal/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrAlreadyLoggedIn            = errors.New("user is already logged in")
	ErrUserNotFound               = errors.New("failed to find user")
	ErrInvalidPassword            = errors.New("invalid password")
	ErrFailedToGenerateSessionKey = errors.New("failed to generate session key")
	ErrFailedToSaveSessionKey     = errors.New("failed to save session key")
	dbColumnUserIdentifier        = "email_address"
)

type UserServiceUtil struct {
	dbManager      db.DBManagerInterface
	sessionManager *session.LoginSessionManager
}

func NewUserServiceUtil() *UserServiceUtil {
	return &UserServiceUtil{
		dbManager:      db.GetDBManager(),
		sessionManager: session.GetLoginSessionManager(),
	}
}

func (usu *UserServiceUtil) AuthenticateUser(loginInfo models.LoginInfo, userSessionKey string) (models.User, error) {
	_, isLoggedIn := usu.checkUserLoggedIn(userSessionKey, loginInfo)

	var user models.User
	err := usu.dbManager.Read(&user, dbColumnUserIdentifier, loginInfo.EmailAddress)
	if err != nil {
		return models.User{}, ErrUserNotFound
	} else if isLoggedIn {
		return user, ErrAlreadyLoggedIn
	}

	if !password.CheckPasswordHash(loginInfo.Password, user.Password) {
		return user, ErrInvalidPassword
	}

	return user, nil
}

func (usu *UserServiceUtil) GenerateSessionKey(user models.User) (string, error) {
	sessionKey, err := session.GenerateRandomSessionKey()
	if err != nil {
		return "", ErrFailedToGenerateSessionKey
	}

	// 세션 키를 캐시에 저장합니다.
	err = usu.sessionManager.SetSession(sessionKey, user.EmailAddress)
	if err != nil {
		return "", ErrFailedToSaveSessionKey
	}

	return sessionKey, nil
}

func (usu *UserServiceUtil) checkUserLoggedIn(userSessionKey string, loginInfo models.LoginInfo) (string, bool) {
	if usu.sessionManager.IsSessionValid(userSessionKey, loginInfo.EmailAddress) {
		fmt.Println("User is already logged in")

		sessionKey, err := usu.sessionManager.GetSession(loginInfo.EmailAddress)
		if err != nil {
			return "", false
		}

		return sessionKey, true
	}
	return "", false
}

// handleLoginError 함수는 로그인 과정에서 발생한 오류를 처리합니다.
// 오류 유형에 따라 적절한 HTTP 상태 코드를 반환합니다.
func (usu *UserServiceUtil) HandleLoginError(ginContext *gin.Context, err error) {
	switch err {
	case ErrAlreadyLoggedIn:
		ginContext.JSON(http.StatusConflict, gin.H{"error": "User is already logged in"})
	case ErrUserNotFound:
		ginContext.JSON(http.StatusNotFound, gin.H{"error": "Failed to find user"})
	case ErrInvalidPassword:
		ginContext.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
	case ErrFailedToGenerateSessionKey, ErrFailedToSaveSessionKey:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process session key"})
	default:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}