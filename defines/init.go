package defines

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type conf struct {
	ResRoot      string;
    LogsRoot     string;
    DmpRoot      string;

    ServerPort   uint;

    RemoteWinSvr string;
    CrashApi     string;
    PdbsApi      string;
}

func createPaths() {
	if err := os.MkdirAll(LogsRoot, 0644); err != nil {
		fmt.Println(err);
	}

	if err := os.MkdirAll(filepath.Join(DmpRoot, "tmp"), 0644); err != nil {
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
	DmpRoot  = c.DmpRoot;

	ServerPort = c.ServerPort;

	RemoteWinSvr = c.RemoteWinSvr;
	CrashApi     = c.CrashApi;
	PdbsApi      = c.PdbsApi;

	createPaths();
}