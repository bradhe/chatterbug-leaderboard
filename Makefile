
clean:
	rm -rf bin
	mkdir bin

build: clean
	GOOS=linux go build -o ./bin/get-chatterbug-leaderboard ./pkg/lambda/get-chatterbug-leaderboard

package: build
	cd ./bin && zip get-chatterbug-leaderboard.zip get-chatterbug-leaderboard
