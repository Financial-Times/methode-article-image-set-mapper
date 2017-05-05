package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
	"os"
	"os/signal"
	"syscall"
	"sync"
)

type app struct {
	args             args
	mapperService    ImageSetMapper
	messageToNativeMapper    MessageToNativeMapper
	queue            queue
	consumerTeardown sync.WaitGroup
	routing          routing
}

func main() {
	cliApp := cli.App("methode-article-image-set-mapper", "Maps inline image-sets from bodies of Methode articles.")
	a := app{}
	a.args = resolveArgs(cliApp)
	cliApp.Action = func() {
		if len(a.args.addresses) == 0 {
			logrus.Fatal("No queue address provided. Quitting...")
		}
		logrus.Infof("methode-article-image-set-mapper is starting systemCode=%s appName=%s port=%s", a.args.appSystemCode, a.args.appName, a.args.port)
		a.messageToNativeMapper = defaultMessageToNativeMapper{}

		a.queue = newQueue(a.args, a.messageToNativeMapper, a.mapperService)
		a.queue.startConsuming()
		a.routing = newRouting(a.messageToNativeMapper, a.mapperService, a.args.appSystemCode, a.args.appName)
		err := a.routing.listenAndServe(a.args.port)
		if err != nil {
			logrus.Fatalf("Cound't serve http endpoints. %v\n", err)
		}
		a.waitForSignals()
		a.teardown()
	}
	err := cliApp.Run(os.Args)
	if err != nil {
		logrus.Fatalf("methode-article-image-set-mapper could not start. %v\n", err)
	}
}

type args struct {
	appSystemCode string
	appName       string
	port          string

	addresses     []string
	group         string
        readTopic     string
        readQueue     string
        writeTopic    string
	writeQueue    string
	authorization string
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

	addresses := app.Strings(cli.StringsOpt{
		Name:   "queue-addresses",
		Desc:   "Addresses to connect to the queue (hostnames).",
		EnvVar: "Q_ADDR",
	})

	group := app.String(cli.StringOpt{
		Name:   "group",
		Desc:   "Group used to read the messages from the queue.",
		EnvVar: "Q_GROUP",
	})

	readTopic := app.String(cli.StringOpt{
		Name:   "read-topic",
		Desc:   "The topic to read the meassages from.",
		EnvVar: "Q_READ_TOPIC",
	})

	readQueue := app.String(cli.StringOpt{
		Name:   "read-queue",
		Desc:   "The queue to read the meassages from.",
		EnvVar: "Q_READ_QUEUE",
	})

	writeTopic := app.String(cli.StringOpt{
		Name:   "write-topic",
		Desc:   "The topic to write the meassages to.",
		EnvVar: "Q_WRITE_TOPIC",
	})

	writeQueue := app.String(cli.StringOpt{
		Name:   "write-queue",
		Desc:   "The queue to write the meassages to.",
		EnvVar: "Q_WRITE_QUEUE",
	})

	authorization := app.String(cli.StringOpt{
		Name:   "authorization",
		Desc:   "Authorization key to access the queue.",
		EnvVar: "Q_AUTHORIZATION",
	})
	return args{
		appSystemCode: *appSystemCode,
		appName:       *appName,
		port:          *port,
		addresses:     *addresses,
		group:         *group,
		readTopic:     *readTopic,
		readQueue:     *readQueue,
		writeTopic:    *writeTopic,
		writeQueue:    *writeQueue,
		authorization: *authorization,
	}
}

func (m app) waitForSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}

func (m app) teardown() {
	m.queue.stop()
	m.consumerTeardown.Wait()
}
