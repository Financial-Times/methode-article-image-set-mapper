package main

import (

	"github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := cli.App("methode-article-image-set-mapper", "Maps inline image-sets from bodies of Methode articles.")
	args := resolveArgs(app)
	logrus.Infof("methode-article-image-set-mapper is starting...\n")
	app.Action = func() {
		logrus.Infof("systemCode=%s appName=%s port=%s", args.appSystemCode, args.appName, args.port)
		mapperService := newImageSetMapper()
		routing := newRouting(mapperService, args.appSystemCode, args.appName)
		err := routing.listenAndServe(args.port)
		if err != nil {
			logrus.Fatalf("Cound't serve http endpoints. %v\n", err)
		}
		waitForSignals()
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatalf("methode-article-image-set-mapper could not start. %v\n", err)
	}
}

type args struct {
	appSystemCode string
	appName string
	port string
}

func resolveArgs(app *cli.Cli) args {
	appSystemCode := app.String(cli.StringOpt{
		Name:   "app-system-code",
		Value:  "methode-article-image-set-mapper",
		Desc:   "System Code of the application",
		EnvVar: "APP_SYSTEM_CODE",
	})
	appName := app.String(cli.StringOpt{
		Name:   "app-name",
		Value:  "methode-article-image-set-mapper",
		Desc:   "Application name",
		EnvVar: "APP_NAME",
	})
	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_PORT",
	})
	return args {
		appSystemCode: *appSystemCode,
		appName: *appName,
		port: *port,
	}
}

func waitForSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
