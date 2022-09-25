package predictions

import (
	"lah-2022/backend/utils"
	"net/http"
	"os"

	cohere "github.com/cohere-ai/cohere-go"
)

type CohereResult struct {
	Prediction string
}

func getPrediction(email, modelId, tokenEnv string) (*CohereResult, *utils.ErrorResponse) {
	token := os.Getenv(tokenEnv)
	co, err := cohere.CreateClient(token)
	if err != nil {
		println(err.Error())
		return nil, &utils.ErrorResponse{Code: http.StatusInternalServerError, Body: utils.APIError{
			ErrorMessage: "cohere problems. 1",
		}}
	}
	response, err := co.Classify(cohere.ClassifyOptions{
		Model:  modelId,
		Inputs: []string{email},
	})

	if err != nil {
		println(err.Error())
		return nil, &utils.ErrorResponse{Code: http.StatusInternalServerError, Body: utils.APIError{
			ErrorMessage: "cohere problems. 2",
		}}
	}

	// The cohere classify model returns a string
	// The proper label is response.classifications[0].prediction
	// proper label is a string
	// 0 is 1
	return &CohereResult{
		Prediction: response.Classifications[0].Prediction,
	}, nil
}
