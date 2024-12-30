package netex

import "net"

// GetFreePort 获取一个未使用的随机的端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	cli, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer cli.Close()
	return cli.Addr().(*net.TCPAddr).Port, nil
}
