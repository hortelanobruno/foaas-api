package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hortelanobruno/foaas-api/constants"
	"github.com/hortelanobruno/foaas-api/ratelimiter"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RateLimiter(rateLimiter ratelimiter.RateLimiter) gin.HandlerFunc {

	return func(c *gin.Context) {
		userID := c.GetHeader(constants.UserIDHeader)

		if !rateLimiter.AllowRequest(userID) {
			logrus.Errorf("Too Many Requests for userID: %s", userID)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": http.StatusText(http.StatusTooManyRequests),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
