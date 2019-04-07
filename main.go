package main

import "crash_server_gate/common"
import "crash_server_gate/defines"
import "crash_server_gate/server"

func main() {
	defines.InitDefines("E:/code/go/src/crash_server_gate/configure/example.json");
	common.InitLogger();
	server.Run();
}
