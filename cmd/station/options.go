package main

type Options struct {
	Host     string
	ClientID string
}

func DefaultOptions() *Options {
	return &Options{
		Host:     "127.0.0.1:1883",
		ClientID: "test-client",
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
