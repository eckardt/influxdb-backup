package main

import (
  "flag"
  "fmt"
  "io"
  "log"
  "os"
  "github.com/eckardt/influxdb-go"
)

type ClientConfig struct {
  *influxdb.ClientConfig
  Destination string
  Series string
}

type Client struct {
  *influxdb.Client
  *ClientConfig
}

func parseFlags() (*ClientConfig) {
  config := &ClientConfig{&influxdb.ClientConfig{}, ""}
  flag.StringVar(&config.Host, "host", "localhost:8086", "host to connect to")
  flag.StringVar(&config.Username, "username", "root", "username to authenticate as")
  flag.StringVar(&config.Password, "password", "root", "password to authenticate with")
  flag.StringVar(&config.Database, "database", "", "database to dump")
  flag.StringVar(&config.Series, "series", "/.*/", "series to dump")
  flag.StringVar(&config.Destination, "out", "-", "output file (default to stdout)")
  flag.BoolVar(&config.IsSecure, "https", false, "connect via https")
  flag.Parse()
  if config.Database == "" {
    fmt.Fprintln(os.Stderr, "flag is mandatory but not provided: -database")
    flag.Usage()
    os.Exit(1)
  }
  if config.Series == "" {
    config.Series = "/.*/"
  }
  return config
}

func main() {
  config := parseFlags()
  _client, err := influxdb.NewClient(config.ClientConfig)
  if err != nil {
    log.Fatal(err)
  }
  client := Client{_client, config}
  client.DumpSeries()
}

func (self *Client) DumpSeries() {
  var err error
  var file io.Writer
  if self.Destination != "-" {
    file, err = os.Create(self.Destination)
    if err != nil {
      log.Fatal(err)
    }
  } else {
    file = os.Stdout
  }
  err = self.QueryStream("SELECT * FROM " + self.Series, file)
  if err != nil {
    log.Fatal(err)
  }
}
