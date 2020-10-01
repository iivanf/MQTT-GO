package main

import (
	"encoding/json"
	"math/rand"
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
	cfgFile, topic string
	host, clientID string
	optsCfg        *Options
	opts           *mqtt.ClientOptions
)

type Sensor struct {
	Name  string  `json:"name"`
	Data  float32 `json:data`
	Units string  `json:"units"`
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "deploy/boat/config/config.toml", "config file (YAML or TOML)")
	rootCmd.Execute()
}

func start() error {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	for {
		sensor := new(Sensor)
		data := rand.Float32()
		sensor.Name = "Sensor-test"
		sensor.Data = data
		sensor.Units = "ºC"
		sensorJSON, err := json.Marshal(sensor)
		if err != nil {
			log.Error(err.Error())
		}

		publicar(sensorJSON)
		time.Sleep(5 * time.Second)
		select {
		case <-stop:
			log.Info("Disconnecting ...")
			// Disconnect
			c.Disconnect(250)
			return nil
		default:

		}
	}
}

func publicar(data []byte) error {
	// Publish a message
	token := c.Publish(optsCfg.Topic, 0, false, data)
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
	topic = viper.GetString("mqtt.topic")
}

func initOpts() {
	optsCfg = DefaultOptions()

	optsCfg.WithHost(host)
	optsCfg.WithClientID(clientID)
	optsCfg.WithTopic(topic)
}

func initClient() {
	opts = mqtt.NewClientOptions().AddBroker(optsCfg.Host).SetClientID(optsCfg.ClientID)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c = mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Error("AQUI %s", token.Error())
	}
}
