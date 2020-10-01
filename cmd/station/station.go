package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	c              mqtt.Client
	cfgFile        string
	host, clientID string
	optsCfg        *Options
	opts           *mqtt.ClientOptions
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Infof("TOPIC: %s\n", msg.Topic())
	log.Infof("MSG: %s\n", msg.Payload())
}

func main() {
	cobra.OnInitialize(initConfig)
	var rootCmd = &cobra.Command{
		Use:   "Mqtt client on station",
		Short: "launch mqqt client for read data of boat",
		Long:  `launch mqqt client for read data of boat`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			initOpts()
			initClient()
			start()
		},
	}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "deploy/station/config/config.toml", "config file (YAML or TOML)")
	rootCmd.Execute()
}

func start() error {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	for {
		if token := c.Subscribe("/boat", 0, nil); token.Wait() && token.Error() != nil {
			log.Error(token.Error())
		}
		select {
		case <-stop:
			log.Info("Disconnecting ...")
			// Unscribe
			if token := c.Unsubscribe("/boat"); token.Wait() && token.Error() != nil {
				log.Error(token.Error())
				return token.Error()
			}
			// Disconnect
			c.Disconnect(250)
			return nil
		default:

		}
	}
}

func publicar(data string) error {
	// Publish a message
	token := c.Publish("/boat", 0, false, data)
	token.Wait()
	return nil
}

func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Errorf("[Init] Unable to read config from file %s: %s", cfgFile, err.Error())
			os.Exit(1)
		} else {
			log.Infof("[Init] Read configuration from file %s", cfgFile)
		}
	}

	host = viper.GetString("mqtt.host")
	clientID = viper.GetString("mqtt.client")
}

func initOpts() {
	optsCfg = DefaultOptions()

	optsCfg.WithHost(host)
	optsCfg.WithClientID(clientID)
}

func initClient() {
	opts = mqtt.NewClientOptions().AddBroker(optsCfg.Host).SetClientID(optsCfg.ClientID)
	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c = mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	}
}
