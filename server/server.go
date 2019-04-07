package server

import "bytes"
import "fmt"
import "io"
import "io/ioutil"
import "net/http"
import "os"
import "path/filepath"
import "strconv"
import "strings"

import "utilities"

import "crash_server_gate/common"
import "crash_server_gate/defines"

var handler *utilities.RequestHandler;
var winSvrClient *http.Client;

func handleDefaultPage(w http.ResponseWriter, r *http.Request) bool {
  if strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path += "index.html";
		http.DefaultServeMux.ServeHTTP(w, r);
		return true;
	}

	return false;
}

func defaultHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("require ", r.URL.Path);
	if handleDefaultPage(w, r) {
		return;
	}

	path := filepath.Join(defines.ResRoot, r.URL.Path[1:]);
	f, err := os.Open(path);
	if err != nil {
		fmt.Println(err);
		w.WriteHeader(404);
		w.Write([]byte(""));
		return;
	}

	data, err2 := ioutil.ReadAll(f);
	if err2 != nil {
		fmt.Println("read err");
		fmt.Println(err2);
		w.WriteHeader(404);
		w.Write([]byte(""));
		return;
	}
	
	w.Write(data);
}

func receiveDmp(w http.ResponseWriter, r *http.Request) {
	buf :=  new(bytes.Buffer);
	if _, err := io.Copy(buf, r.Body); err != nil {
		common.ErrorLogger.Print(err);
		return;
	}

	nr, err2 := http.NewRequest(http.MethodPost, defines.RemoteWinSvr + defines.CrashApi, buf);
	if err2 != nil {
		common.ErrorLogger.Print(err2);
		return;
	}

	nr.Header.Add("Content-Type", r.Header.Get("Content-Type"));

	go RedirtDmp(nr);
}

func receivePdbs(w http.ResponseWriter, r *http.Request) {
	buf :=  new(bytes.Buffer);
	if _, err := io.Copy(buf, r.Body); err != nil {
		common.ErrorLogger.Print(err);
		return;
	}

	nr, err2 := http.NewRequest(http.MethodPost, defines.RemoteWinSvr + defines.PdbsApi, buf);
	if err2 != nil {
		common.ErrorLogger.Print(err2);
		return;
	}

	nr.Header.Add("Content-Type", r.Header.Get("Content-Type"));

	go RedirtDmp(nr);
}

func RedirtDmp(r *http.Request) {
	res, err := winSvrClient.Do(r);
	fmt.Println(res);
	if err != nil {
		common.ErrorLogger.Print(err);
		return;
	}

	if res.StatusCode != http.StatusOK {
		common.ErrorLogger.Print(res.StatusCode);
		return;
	}
}

func initServer() {
	handles := make(map[string]func(http.ResponseWriter, *http.Request));
	handles["/postCrash"] = receiveDmp;
	handles["/postPdbs"] = receivePdbs;
	handles["/"] = defaultHandle;

	handler = utilities.NewRequestHandler(&handles);

	winSvrClient = new(http.Client);
}

func Run() {
	initServer();

	p := ":" + strconv.FormatUint(uint64(defines.ServerPort), 10);
	http.ListenAndServe(p, handler);
}