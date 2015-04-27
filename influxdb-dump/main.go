package main

import (
  "flag"
  "fmt"
  "io"
  "log"
  "os"
  "strings"
  "github.com/eckardt/influxdb-go"
)

type ClientConfig struct {
  *influxdb.ClientConfig
  Destination string
  Series string
  StartTime string
  EndTime string
}

type Client struct {
  *influxdb.Client
  *ClientConfig
}

func parseFlags() (*ClientConfig) {
  config := &ClientConfig{&influxdb.ClientConfig{}, "", "", "", ""}
  flag.StringVar(&config.Host, "host", "localhost:8086", "host to connect to")
  flag.StringVar(&config.Username, "username", "root", "username to authenticate as")
  flag.StringVar(&config.Password, "password", "root", "password to authenticate with")
  flag.StringVar(&config.Database, "database", "", "database to dump")
  flag.StringVar(&config.Series, "series", "/.*/", "series name to dump")
  flag.StringVar(&config.StartTime, "start", "", "time since series must be dumped")
  flag.StringVar(&config.EndTime, "end", "", "time till series must be dumped")
  flag.StringVar(&config.Destination, "out", "-", "output file (default to stdout)")
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
  client.DumpSeries()
}

func (self *Client) DumpSeries() {
  var err error
  var file io.Writer
  var query string
  var where_time string

  times := make([]string, 0)

  if self.Destination != "-" {
    file, err = os.Create(self.Destination)
    if err != nil {
      log.Fatal(err)
    }
  } else {
    file = os.Stdout
  }

  if self.StartTime != "" {
    times = append(times, fmt.Sprintf("time > %s", self.StartTime))
  } 

  if self.EndTime != "" {
    times = append(times, fmt.Sprintf("time < %s", self.EndTime))
  }

  if len(times) > 0 {
    where_time = fmt.Sprintf("where %s", strings.Join(times, " and "))
  }else{
    where_time = ""
  }

  query = fmt.Sprintf("SELECT * FROM %s %s", self.Series, where_time)

  log.Println("Query:", query)

  err = self.QueryStream(query, file)
  if err != nil {
    log.Fatal(err)
  }
}
