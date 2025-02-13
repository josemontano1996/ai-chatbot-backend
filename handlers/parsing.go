package handlers

import (

	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/sharedtypes"
)

func ParseUserMessageFromRequest(c *gin.Context) (*sharedtypes.Message, error) {

	var userMsg sharedtypes.Message

	if err := c.ShouldBindJSON(&userMsg); err != nil {
		return nil, err
	}

	userMsg.Code = sharedtypes.UserMessageCode

	return &userMsg, nil
}
