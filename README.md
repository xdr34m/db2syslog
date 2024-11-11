# db2syslog
Golang CLI that reads DB Table and forwards rows to syslog <br>
read&write needed, so it can mark processed rows in the table

# supported DBs
mysql, oracle

# config example
```yaml
# Database configuration
database:
  type: "mysql"  # Options: "mysql" or "oracle"
  host: "localhost"
  port: 3306
  username: "user"
  password: "password"
  dbname: "logsdb"
  table: "logs"
  keycolumn: "id"
  markedcolumn: "processed"
  unmarkedvalue: "false"
  markedvalue: "true"
  presets: # only for oracle - this translates to ALTER SESSION SET <yamlkey> = '<yamlvalue>'
    #time_zone: "Europe/Berlin"
    nls_date_format: "DD.MM.RR HH24:MI:SS"
    nls_timestamp_format: "DD.MM.YYYY HH24:MI:SS.FF"
    nls_timestamp_tz_format: "DD.MM.RR HH24:MI:SSXFF"
    #nls_language: "GERMAN"
    #nls_territory: "GERMANY"

# Syslog configuration
syslog:
  server: "localhost"
  port: 10312
  transport: "tcp" # Options: "tcp" or "udp"
  hostname: "myhost"
  app_name: "myapp"
  proc_id: "12345"
  msg_id: "-"  # Optional, can be left as "-"
  # Structured Data (SDATA) example is needed...
  sdata: "exampleSDID@0815 eventSource=\"Application\" eventID=\"0815\""
  dropcolumn:
    - id
    - column1

# Optional logging level or other configurations
log_level: "info"
