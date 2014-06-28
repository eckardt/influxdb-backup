package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "io"
  "log"
  "os"
  "github.com/eckardt/influxdb-go"
)

type ClientConfig struct {
  *influxdb.ClientConfig
  Source string
  Files []string
}

type Client struct {
  *influxdb.Client
  *ClientConfig
}

func parseFlags() (*ClientConfig) {
  config := &ClientConfig{&influxdb.ClientConfig{}, "", nil}
  flag.StringVar(&config.Host, "host", "localhost:8086", "host to connect to")
  flag.StringVar(&config.Username, "username", "root", "username to authenticate as")
  flag.StringVar(&config.Password, "password", "root", "password to authenticate with")
  flag.StringVar(&config.Database, "database", "", "database to restore")
  flag.StringVar(&config.Source, "in", "-", "input file (default stdin)")
  flag.BoolVar(&config.IsSecure, "https", false, "connect via https")
  flag.Parse()
  if config.Database == "" {
    fmt.Fprintln(os.Stderr, "flag is mandatory but not provided: -database")
    flag.Usage()
    os.Exit(1)
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
  client.DisableCompression()
  client.ImportSeries()
}

func (self *Client) ImportSeries() {
  var err error
  var file io.Reader
  if self.Source != "-" {
    file, err = os.Open(self.Source)
    if err != nil {
      log.Fatal(err)
    }
  } else {
    file = os.Stdin
  }

  dec := json.NewDecoder(file)
  for {
    series := make([]*influxdb.Series, 1)
    if err := dec.Decode(&series[0]); err == io.EOF {
      break
    } else if err != nil {
      log.Fatal(err)
    }
    err := self.WriteSeries(series)
    if err != nil {
      log.Fatal(err)
    }
  }
}
