package types

type HostInfo struct {
  ID int
  Name string
  User string
  Port int
  Online bool
  PublicIP string
  LocalIP string
  LastOnline string
  Version int
}

type HostConfig struct {
  Server ConnectionInfo
  SSHStartStop bool
  ID int
  Name string
  User string
  Port int
}

type HostInfoConfig struct {
  Verbose bool
  IDLength int
  NameLength int
  UserLength int
  PortLength int
  PublicIPLength int
  LocalIPLength int
  LastOnlineLength int
}

type ServerInfo struct {
  Name string
  DbPort string
  SshPort string
  WebPort string
  UserPassword string
  RootPassword string
  WebInterface bool
}

type ConnectionInfo struct {
  Name string
  Host string
  SshPort string
  DbPort string
  Password string
}

type ClientConfig struct {
  ConnectionConfigs []ConnectionInfo
  HostConfigs []HostConfig
}

type ReverseTunnelConfig struct {
  ConnectionInfo ConnectionInfo
  LocalPort int
  RemotePort int
}

type ForwardTunnelConfig struct {
  ConnectionInfo ConnectionInfo
  LocalPort int
  RemotePort int
}
