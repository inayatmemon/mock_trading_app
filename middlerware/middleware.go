package middleware

import (
	commonmodels "9DotTechnology/trading_platform/models/common"
	"9DotTechnology/trading_platform/utils/auth"
	"9DotTechnology/trading_platform/utils/commons"

	"9DotTechnology/trading_platform/utils/localization"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LanHeader - is for headers without auth
type LanHeader struct {
	Language    string `header:"language" binding:"required"`
	ContentType string `header:"content-type" binding:"required"`
}

// AuthLanHeader - is for headers with auth
type AuthLanHeader struct {
	Authorization string `header:"authorization" binding:"required"`
	Language      string `header:"language" binding:"required"`
	ContentType   string `header:"content-type" binding:"required"`
}

// This middleware sets whether the user is logged in or not
func ValidateHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := LanHeader{}
		err := c.BindHeader(&headers)
		if err != nil {
			commons.SendApiResponse(c, commonmodels.ApiResponse{
				Code:    http.StatusBadRequest,
				Message: "language and content-type headers are required.",
			})
			return
		}
	}
}

// This middleware sets whether the user is logged in or not
func ValidateAuthHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := AuthLanHeader{}
		err := c.BindHeader(&headers)
		if err != nil {
			commons.SendApiResponse(c, commonmodels.ApiResponse{
				Code:    http.StatusBadRequest,
				Message: "authorization, language and content-type headers are required.",
			})
			return
		}
	}
}

func ValidateAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("language")
		authToken := c.GetHeader("authorization")
		if authToken == "" {
			commons.SendApiResponse(c, commonmodels.ApiResponse{
				Code:    http.StatusBadRequest,
				Message: "authorization header is required.",
			})
			return
		}

		_, err := auth.DecodeToken(authToken)
		if err != nil {
			logwrapper.Logger.Debugln("error in decoding authToken: ", err)
			commons.SendApiResponse(c, commonmodels.ApiResponse{
				Code:    http.StatusUnauthorized,
				Message: localization.GetMessage(lang, "common.401", nil),
			})
			return
		}
	}
}

func ValidateStreamConnectionAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for a specific header (e.g., "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header is required.", http.StatusBadRequest)
			logwrapper.Logger.Debugln("Authorization header missing")
			return
		}

		_, err := auth.DecodeToken(authHeader)
		if err != nil {
			http.Error(w, localization.GetMessage("en", "common.401", nil), http.StatusUnauthorized)
			logwrapper.Logger.Debugln("Authorization header invalid")
			return
		}

		// If all validations pass, proceed to the next handler
		next.ServeHTTP(w, r)
	}
}
