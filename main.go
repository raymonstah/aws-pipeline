package main

import (
	"log"
	"os"

	"github.com/raymonstah/aws-pipeline/internal/awsutils"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Usage = "deploy app to AWS"
	app.EnableBashCompletion = true
	app.HideVersion = true

	app.Commands = []*cli.Command{
		awsutils.LambdaCommand,
		awsutils.CloudformationCommand,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
