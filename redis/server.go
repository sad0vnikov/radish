package redis

//Server is a struct storing redis server parameters
type Server struct {
	Name string
	Host string
	Port int
}

//NewServer returns a redis.Server struct with given fields
func NewServer(name, host string, port int) Server {
	return Server{Name: name, Host: host, Port: port}
}
