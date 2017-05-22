package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/Financial-Times/message-queue-go-producer/producer"
	"net/http"
	"net"
	"time"
	"fmt"
)

type app struct {
	args             args
	queue            defaultQueue
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	cliApp := cli.App("methode-article-image-set-mapper", "Maps inline image-sets from bodies of Methode articles.")
	a := app{}
	a.args = resolveArgs(cliApp)
	cliApp.Action = func() {
		if len(a.args.addresses) == 0 {
			logrus.Fatal("No queue address provided. Quitting...")
		}
		logrus.Infof("methode-article-image-set-mapper is starting systemCode=%s appName=%s port=%s", a.args.appSystemCode, a.args.appName, a.args.port)
		messageToNativeMapper := defaultMessageToNativeMapper{}
		imageSetMapper := newImageSetMapper(defaultArticleToImageSetMapper{}, defaultAttributesMapper{}, defaultImageSetToJSONMapper{})
		httpClient := http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConnsPerHost:   20,
				TLSHandshakeTimeout:   3 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}
		consumerConfig := consumer.QueueConfig{
			Addrs:                a.args.addresses,
			Group:                a.args.group,
			Topic:                a.args.readTopic,
			Queue:                a.args.readQueue,
			ConcurrentProcessing: false,
			AutoCommitEnable:     true,
			AuthorizationKey:     a.args.authorization,
		}
		producerConfig := producer.MessageProducerConfig{
			Addr:          a.args.addresses[0],
			Topic:         a.args.writeTopic,
			Queue:         a.args.writeQueue,
			Authorization: a.args.authorization,
		}
		prettyPrintConfig(consumerConfig, producerConfig)
		messageProducer := producer.NewMessageProducerWithHTTPClient(producerConfig, &httpClient)
		a.queue = newQueue(nil, messageProducer, messageToNativeMapper, imageSetMapper)
		messageConsumer := consumer.NewConsumer(consumerConfig, a.queue.onMessage, &httpClient)
		a.queue.messageConsumer = messageConsumer
		a.queue.startConsuming()
		httpMappingHandler := newHTTPMappingHandler(messageToNativeMapper, imageSetMapper)
		routing := newRouting(httpMappingHandler, &httpClient, consumerConfig, a.args.appSystemCode, a.args.appName)
		go routing.listenAndServe(a.args.port)
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
	var consumerTeardown sync.WaitGroup
	m.queue.stop(consumerTeardown)
	consumerTeardown.Wait()
	logrus.Infof("Doneee$$$")
}

func prettyPrintConfig(consumerConfig consumer.QueueConfig, producerConfig producer.MessageProducerConfig) string {
	return fmt.Sprintf("Config: [\n\t%s\n\t%s\n]", prettyPrintConsumerConfig(consumerConfig), prettyPrintProducerConfig(producerConfig))
}

func prettyPrintConsumerConfig(consumerConfig consumer.QueueConfig) string {
	return fmt.Sprintf("consumerConfig: [\n\t\taddr: [%v]\n\t\tgroup: [%v]\n\t\ttopic: [%v]\n\t\treadQueueHeader: [%v]\n\t]",
		consumerConfig.Addrs, consumerConfig.Group, consumerConfig.Topic, consumerConfig.Queue)
}

func prettyPrintProducerConfig(producerConfig producer.MessageProducerConfig) string {
	return fmt.Sprintf("producerConfig: [\n\t\taddr: [%v]\n\t\ttopic: [%v]\n\t\twriteQueueHeader: [%v]\n\t]",
		producerConfig.Addr, producerConfig.Topic, producerConfig.Queue)
}
