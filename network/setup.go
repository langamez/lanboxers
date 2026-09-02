package network

func SetupConnection(isHost bool, ip string, port int) (*Connection, error) {
	if isHost {
		return StartServer(port)
	}

	return Dial(ip, port)
}
