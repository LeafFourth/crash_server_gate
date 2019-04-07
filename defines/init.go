package defines

import (
	"encoding/json"
	"fmt"
	"os"
)

type conf struct {
	ResRoot      string;
    LogsRoot     string;

    ServerPort   uint;

    RemoteWinSvr string;
    CrashApi     string;
    PdbsApi      string;
}

func createPaths() {
	err := os.MkdirAll(LogsRoot, 0644);

	if err != nil {
		fmt.Println(err);
	}
}

func InitDefines(confFile string) {
	r, err := os.Open(confFile);
	if err != nil {
		fmt.Println("1:", err);
		return;
	}

	d := json.NewDecoder(r);
	if d == nil {
		fmt.Println("2", "out of memory");
		return;
	}

	var c conf;
	d.Decode(&c);

	ResRoot  = c.ResRoot;
	LogsRoot = c.LogsRoot;

	ServerPort = c.ServerPort;

	RemoteWinSvr = c.RemoteWinSvr;
	CrashApi     = c.CrashApi;
	PdbsApi      = c.PdbsApi;

	createPaths();

}