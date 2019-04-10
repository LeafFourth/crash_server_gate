package server

import "bytes"
import "encoding/json"
import "fmt"
import "io"
import "io/ioutil"
import "net/http"
import "os"
import "path/filepath"
import "strconv"
import "strings"
import "time"

import "utilities"

import "crash_server_gate/common"
import "crash_server_gate/defines"
import "crash_server_gate/db"

type dmpDesc struct {
	Uid  int;
	Ver  string;
	Date string;
	Name string;
}

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
	//fmt.Println(nr.Header.Get("Content-Type"));
	//fmt.Println(nr);
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
	if err != nil {
		common.ErrorLogger.Print(err);
		return;
	}

	if res.StatusCode != http.StatusOK {
		common.ErrorLogger.Print(res.StatusCode);
		return;
	}
}

func formatDate(date string) time.Time {
	length := len(date);
	if length < 5 {
		return time.Now();
	}

	d, _ := strconv.ParseInt(date[length - 2:], 10, 0);
	m, _ := strconv.ParseInt(date[length - 4: length - 2], 10, 0);
	y, _ := strconv.ParseInt(date[:length - 4], 10, 0);

	if d == 0 || m == 0 {
		return time.Now();
	}

	return time.Date(int(y), time.Month(m), int(d), 0, 0, 0, 0, time.UTC);
}

func writeDb(ds dmpDesc, cs string) {
	if ds.Uid == 0 {
		common.ErrorLogger.Print("uid excepted!");
		return;
	}

	date := formatDate(ds.Date);

	db.PreCreateTableForDate(date);

	c := db.GetConn();
	if c == nil {
		common.ErrorLogger.Print("db nil");
		return;
	}

	tableName := db.GetTableName(date);

	if len(cs) > defines.CSFieldLen {
		cs = cs[:defines.CSFieldLen];
	}

	_, err := c.Exec("INSERT INTO " + tableName + "(name, uid, callstack) VALUES(?, ?, ?)", ds.Name, ds.Uid, cs);
	if err != nil {
		common.ErrorLogger.Print(err);
		return;
	}
}

func ReceiveCallStack(w http.ResponseWriter, r *http.Request) {
	cs := r.FormValue("callback");

	ei := r.FormValue("einfo");

	buf := bytes.NewBufferString(ei);
	var ds dmpDesc;
	d := json.NewDecoder(buf);
	err := d.Decode(&ds);
	if err != nil {
		common.ErrorLogger.Print(err);
		return;
	}

	writeDb(ds, cs);
}

func initServer() {
	handles := make(map[string]func(http.ResponseWriter, *http.Request));
	handles["/postCrash"] = receiveDmp;
	handles["/postPdbs"] = receivePdbs;
	handles["/RecvCallstack"] = ReceiveCallStack;
	handles["/"] = defaultHandle;

	handler = utilities.NewRequestHandler(&handles);

	winSvrClient = new(http.Client);
}

func Run() {
	initServer();

	p := ":" + strconv.FormatUint(uint64(defines.ServerPort), 10);
	http.ListenAndServe(p, handler);
}