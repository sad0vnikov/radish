package redis

//Server is a struct storing redis server parameters
type Server struct {
	Name                  string
	Host                  string
	Port                  int
	DatabasesCount        uint8
	ConnectionCheckPassed bool
	KeyspaceStat          map[string]ServerKeyspaceStat
	ServerStat            ServerStat
}

type ServerStat struct {
	ConnectedClientsCount int64
	RedisVersion          string
	UptimeInSeconds       int64
	UsedMemoryHuman       string
	UsedMemoryBytes       int64
	MaxMemoryHuman        string
	MaxMemoryBytes        int64
}
type ServerKeyspaceStat struct {
	KeysCount int64
}

//NewServer returns a redis.Server struct with given fields
func NewServer(name, host string, port int) Server {
	return Server{Name: name, Host: host, Port: port}
}
