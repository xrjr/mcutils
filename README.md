# mcutils - Implementation of Minecraft protocols in Go



<p align="center">
 <img alt="logo" src="https://i.imgur.com/nIBQcRj.png" />
</p>



## Informations

![test workflow](https://github.com/xrjr/mcutils/actions/workflows/tests.yml/badge.svg)

### General

All protocols are implemented in Go, without any external dependency. All protocols should be supported on any platform/architecture as long as Go can compile them.

All protocols have been implemented using [wiki.vg](https://wiki.vg/). All of them are 100% compliant with the standard described there.

This project also contains an helper in communication with tcp/udp, called `pkg/networking`.

This project has no dependency.

### Supported protocols

- [Ping](https://wiki.vg/Server_List_Ping) (Server List Ping)

- [Query](https://wiki.vg/Query)

- [Rcon](https://wiki.vg/Rcon)

  

> All protocols implementations support SRV record resolving. 
>
> Rcon implementation supports fragmented response packets.
>
> Ping protocol changed in 1.7 in a non-backwards compatible way. Both modern and legacy SLP protocols are supported.



## CLI

### Install

```shell
go install github.com/xrjr/mcutils/cmd/mcutils@latest
```

### Usage

```shell
$ mcutils ping <hostname> <port>
Example : mcutils ping localhost 25565

$ mcutils ping-legacy <hostname> <port>
Example : mcutils ping localhost 25565

$ mcutils ping-legacy-1.6.4 <hostname> <port>
Same as ping-legacy, but uses 1.6+ SLP protocol
Example : mcutils ping localhost 25565

$ mcutils query <basic|full> <hostname> <port>
Example : mcutils query basic localhost 25565

$ mcutils rcon <hostname> <port> <password> <command>
Example : mcutils rcon localhost 25575 mypassword "say hello"
```



## How to use (simple way) ?

### Ping

```go
// Ping returns the server list ping infos (JSON-like object), and latency of a minecraft server.
properties, latency, err := ping.Ping("localhost", 25565)
```

```go
// Ping returns the legacy server list ping infos, and latency of a minecraft server.
properties, latency, err := ping.PingLegacy("localhost", 25565)
```

```go
// Ping returns the legacy server list ping infos (using 1.6+ SLP protocol), and latency of a minecraft server.
properties, latency, err := ping.PingLegacy1_6_4("localhost", 25565)
```

### Query

```go
// QueryBasic returns the basic stat of a minecraft server.
basicStat, err := query.QueryBasic("localhost", 25565)

// QueryBasic returns the full stat of a minecraft server.
fullStat, err := query.QueryFull("localhost", 25565)
```

### Rcon

```go
// Rcon executes a command on a minecraft server, and returns the response of that command.
response, err := rcon.Rcon("localhost", 25575, "password", "command")
```



## How to use (full control way) ?

### Ping

```go
pingclient := ping.NewClient("localhost", 25565)

// Connect opens the connection, and can raise an error for example if the server is unreachable
err := pingclient.Connect()

// Handshake is the base request of ping, the one that displays number of players, MOTD, etc...
// If all went well, hs contains a field Properties which contains a golang-usable JSON Object
hs, err := pingclient.Handshake()

// Ping is a request that basically do nothing and is just used for measuring the latency
// pong contains the latency in ms
pong, err := pingclient.Ping()

// Disconnect closes the connection
err = pingclient.Disconnect()
```

```go
legacypingclient := ping.NewClientLegacy("localhost", 25565)

// Connect opens the connection, and can raise an error for example if the server is unreachable
err := legacypingclient.Connect()

// Ping requests for informations and latency of a server, using legacy SLP.
infos, latency, err := legacypingclient.Ping()

// Same as Ping, but uses 1.6+ SLP protocol. Only one call to Ping or Ping1_6_4 should be done per opened connection.
infos, latency, err := legacypingclient.Ping1_6_4()

// Disconnect closes the connection
err = legacypingclient.Disconnect()
```

### Query

```go
queryclient := query.NewClient("localhost", 25565)

// Connect opens the connection, and can raise an error for example if the server is unreachable
err := queryclient.Connect()

// Handshake request is used to get the challenge token, needed for questing basic and full stat
challengeToken, err := queryclient.Handshake()

// BasicStat returns several informations about the server like number of players, maximum number of players, etc... in a fully predictable way
bs, err := queryclient.BasicStat(challengeToken)

// FullStat returns several informations (more than BasicStat) in a JSON format, plus the list of connected players
fs, err := queryclient.FullStat(challengeToken)

// Disconnect closes the connection
queryclient.Disconnect()
```

### Rcon

```go
rconclient := rcon.NewClient("localhost", 25575)

// Connect opens the connection, and can raise an error for example if the server is unreachable
err := rconclient.Connect()

// Authenticate request is used to authenticate the connection.
// If the authentication succeeds, ok will be true, and if it fails, ok will be false.
// err will be nil unless there is a communication problem with the server
ok, err := rconclient.Authenticate("password")

// Command will execute the given command on the server, and the output text will be returned in res
res, err := rconclient.Command("playerlist")

// Disconnect closes the connection
rconclient.Disconnect()
```
