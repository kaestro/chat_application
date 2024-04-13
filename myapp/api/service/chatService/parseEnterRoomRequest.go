// myapp/api/service/chatService/parseEnterRoomRequest.go
package chatService

import (
	"myapp/api/models"

	"github.com/gin-gonic/gin"
)

func IsUpgradeHeaderValid(c *gin.Context) bool {
	upgradeHeader := c.GetHeader("Upgrade")
	return upgradeHeader == "websocket"
}

func ParseEnterRoomRequest(c *gin.Context) (models.RoomRequest, error) {
	var req models.RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return models.RoomRequest{}, err
	}
	return req, nil
}