package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bradhe/chatterbug-leaderboard/pkg/chatterbug"
	"github.com/bradhe/chatterbug-leaderboard/pkg/config"
	"github.com/bradhe/chatterbug-leaderboard/pkg/logs"
	"github.com/bradhe/chatterbug-leaderboard/pkg/slack"
)

type APIGatewayResponse struct {
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
}

func HandleRequest(ctx context.Context) (APIGatewayResponse, error) {
	conf := config.New()

	if conf.Debug {
		logs.EnableDebug()
		logger.Debug("debug mode enabled")
	}

	client := chatterbug.New(conf.ChatterbugAPIToken)

	if leaderboard, err := client.GetLeaderboard(); err != nil {
		logger.WithError(err).Errorf("failed to get leaderboard from chatterbug")
		return APIGatewayResponse{}, fmt.Errorf("failed to get leaderboard from chatterbug")
	} else {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("Top 10 Chatterbug users for *%s*", time.Now().Format("January")))
		builder.WriteString("\n")
		builder.WriteString("\n")

		for i, pos := range leaderboard.TopPositions {
			switch i {
			case 0:
				builder.WriteString(fmt.Sprintf(":first_place: *%s* _(%0.1f hours)_\n", pos.User.Login, pos.StudyTime))
			case 1:
				builder.WriteString(fmt.Sprintf(":second_place: *%s* _(%0.1f hours)_\n", pos.User.Login, pos.StudyTime))
			case 2:
				builder.WriteString(fmt.Sprintf(":third_place: *%s* _(%0.1f hours)_\n", pos.User.Login, pos.StudyTime))
			default:
				builder.WriteString(fmt.Sprintf(":white_circle: *%s* _(%0.1f hours)_\n", pos.User.Login, pos.StudyTime))
			}
		}

		builder.WriteString("\n")
		builder.WriteString("\n")
		builder.WriteString("Visit [Chatterbug](https://app.chatterbug.com) to study more!")

		message := slack.NewMessage(builder.String())

		resp := APIGatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: message.ToJSON(),
		}

		return resp, nil
	}
}

func main() {
	lambda.Start(HandleRequest)
}
