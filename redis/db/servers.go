package db

//GetMaxDbNumsForServer returns a maximum value for db name for a given server
func GetMaxDbNumsForServer(serverName string) (uint8, error) {
	return connector.GetMaxDbNumsForServer(serverName)
}
