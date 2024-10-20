package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rebeccaskinner/gofpher/either"
	"github.com/sirupsen/logrus"

	"github.com/rebeccaskinner/agile17-sample/user"
)

func main() {
	config, err := getArgs()
	if err != nil {
		fmt.Println(err)
		fmt.Println(showHelp())
		os.Exit(1)
	}

	var (
		getEndpoint  = fmt.Sprintf("%s/oldusers/%s", config.endpoint, config.username)
		postEndpoint = fmt.Sprintf("%s/newusers/%s", config.endpoint, config.username)
		get          = either.WrapEither(http.Get)
		body         = func(r *http.Response) io.Reader { return r.Body }
		read         = either.WrapEither(io.ReadAll)
		fromjson     = either.WrapEither(user.NewFromJSON)
		mkUser       = either.WrapEither(user.NewUserFromUser)
		toJSON       = either.WrapEither(json.Marshal)
		updateUser   = either.WrapEither(func(b *bytes.Buffer) (*http.Response, error) {
			return http.Post(postEndpoint, "application/json", b)
		})
		printResponse = either.WrapEither(func(r *http.Response) (*http.Response, error) {
			logrus.Info("resp: %s", r)
			return r, nil
		})
	)

	result := get(getEndpoint).
		LiftM(body).
		AndThen(read).
		AndThen(fromjson).
		AndThen(mkUser).
		AndThen(toJSON).
		LiftM(bytes.NewBuffer).
		AndThen(updateUser).
		Next(printResponse)

	// todo: where is error handling ?
	//fmt.Println(result)
	logrus.Infof("result: %+v", result)
}

type config struct {
	endpoint string
	username string
}

func getArgs() (*config, error) {
	args := os.Args
	if len(args) < 3 {
		return nil, errors.New("insufficient number of arguments")
	}
	switch args[1] {
	case "-?", "-h", "--help":
		return nil, errors.New("showing help message")
	}
	return &config{endpoint: args[1], username: args[2]}, nil
}

func showHelp() string {
	return `client: a simple client to get json data about a user
Usage:
client <endpoint> <username>

Client will connect to server at <endpoint> and request data about <username>.

Example:
client http://localhost:8080 user1      # Get information about user1
`
}
