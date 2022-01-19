package main

import (
	"github.com/bradhe/chatterbug-leaderboard/pkg/chatterbug"
	"github.com/bradhe/chatterbug-leaderboard/pkg/config"
	"github.com/bradhe/chatterbug-leaderboard/pkg/logs"
)

func main() {
	conf := config.New()

	if conf.Debug {
		logs.EnableDebug()
		logger.Debug("debug mode enabled")
	}

	client := chatterbug.New(conf.ChatterbugAPIToken)

	leaderboard, _ := client.GetLeaderboard()
	logger.Infof("Leaderboard: %v", leaderboard)
}
