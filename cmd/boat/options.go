package main

type Options struct {
	Host     string
	ClientID string
	Topic    string
}

func DefaultOptions() *Options {
	return &Options{
		Host:     "tcp://127.0.0.1:1883",
		ClientID: "boat-client",
		Topic:    "/boat",
	}
}

func (opts *Options) WithClientID(clientID string) *Options {
	opts.ClientID = clientID
	return opts
}

func (opts *Options) WithHost(host string) *Options {
	opts.Host = host
	return opts
}

func (opts *Options) WithTopic(topic string) *Options {
	opts.Topic = topic
	return opts
}
