package predictions

import (
	"lah-2022/backend/auth"
	"lah-2022/backend/ent"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type PredictionResult struct {
	FromDomain string `json:"from_domain"`
	Sentiment  bool   `json:"sentiment"`
	Spam       bool   `json:"is_spam"`
}

type api struct {
	client *ent.Client
}

func RegisterRoutes(client *ent.Client, g *echo.Group) {
	a := api{client}
	g.Use(auth.AuthMiddleware)
	g.GET("/email/:emailId", a.handlePrediction)
}

func (a api) handlePrediction(cx echo.Context) error {
	c := cx.(*auth.AuthedContext)

	emailId := c.Param("emailId")

	email, err := getEmail(c.Request().Context(), emailId, c.RefreshToken)
	if err != nil {
		return c.JSON(err.Code, err.Body)
	}

	// Spam
	spamPrediction, err := getPrediction(email.Body, "e7100805-b450-4f36-b544-fda8484f6325-ft", "COHERE_TOKEN")
	if err != nil {
		return c.JSON(err.Code, err.Body)
	}

	isSpam := true
	if spamPrediction.Prediction == "0" {
		isSpam = false
	}

	// Sentiment
	sentimentPrediction, err := getPrediction(email.Body, "e12c7d2d-32f6-4480-af8c-908d442a347b-ft", "COHERE_TOKEN_1")
	if err != nil {
		return c.JSON(err.Code, err.Body)
	}

	senderDomain := strings.Split(email.Sender.Address, "@")[1]

	sentiment := true
	if sentimentPrediction.Prediction == "0" {
		sentiment = false
	}

	return c.JSON(http.StatusOK, PredictionResult{
		FromDomain: senderDomain,
		Sentiment:  sentiment,
		Spam:       isSpam,
	})
}
