package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/cli"
)

type Environment struct {
	MonitoredURLs	[]string `required:"true" envconfig:"MONITORED_URLS"`
}

type Notifier interface {
	Notify(message string) bool
}

func main() {
	var notifierParameter string

	err := godotenv.Load("/go/bin/.env")
	if err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	var e Environment
	err = envconfig.Process("checker", &e)
	if err != nil {
		log.Fatalf("envconfig.Process: %w", err)
	}

	app := cli.NewApp()
	app.Name = "checker"
	app.Usage = "Check if URLs are up"

	app.Flags = []cli.Flag {
		&cli.StringFlag{
		  Name:        "notifier, n",
		  Value:       "slack",
		  Usage:       "Chose a notifier. Supported values : slack",
		  Destination: &notifierParameter,
		},
	  }

	app.Commands = []*cli.Command{
		{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "check if URLs are up",
			Action: func(c *cli.Context) error {
				var statusCode int

				notifier := getNotifier(notifierParameter)

				for _, monitoredURL := range e.MonitoredURLs {
					log.Println(fmt.Sprintf("Visiting %v", monitoredURL))
					resp, err := http.Get(monitoredURL)
					if err != nil {
						log.Println(fmt.Sprintf("Could not connect to url %s : %s", monitoredURL, err))
						statusCode = 0
					} else {
						log.Println(fmt.Sprintf("Response for url %s : %d", monitoredURL, resp.StatusCode))
						statusCode = resp.StatusCode
					}

					switch statusCode {
					case 200:
						log.Println(fmt.Sprintf("URL %s is up, nothing to do.", monitoredURL))
					case 0:
						log.Println(fmt.Sprintf("URL %s is not up (connection error), notifying", monitoredURL))
						message := fmt.Sprintf(":warning: Alert, site %s looks down (connection error : %s) ! :warning:", monitoredURL, err)
						notifier.Notify(message)
					default:
						log.Println(fmt.Sprintf("URL %s is not up (code %d), notifying", monitoredURL, statusCode))
						message := fmt.Sprintf(":warning: Alert, site %s looks down (status code %d) ! :warning:", monitoredURL, statusCode)
						notifier.Notify(message)
					}
				}
				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Test notifier settings",
			Action: func(c *cli.Context) error {
				log.Println("Sending test message")

				notifier := getNotifier(notifierParameter)
				notifier.Notify(":warning: Alert, this is a test message for checker !! :warning:")
				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}


func getNotifier(notifierParameter string) (Notifier) {

  switch notifierParameter {
  case "slack":
    return makeSlackNotifier()
  }

  log.Println("No notifier found, or notifier is not valid. Falling back to slack notifier")
  return makeSlackNotifier()
}
