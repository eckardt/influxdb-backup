# InfluxDB backup and restore

Utility tools to backup and restore InfluxDB databases.

`influxdb-dump` dumps all series from an InfluxDB to a file.

`influxdb-restore` writes all series from a file to an InfluxDB.

## Usage

To copy all datapoints (all series) from one database to another do:

```sh
$ influxdb-dump -database oldDB | influxdb-restore -database newDB
```

See `influxdb-dump -help` for more usage information.

License MIT
