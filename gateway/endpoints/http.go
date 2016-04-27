package endpoints

import (
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	textTemplate "html/template"
	"time"
)

var form = `
<!doctype html>
<html>
	<head>
		<title>EASY LAB Gateway Configuration</title>
	</head>
	<body>
		<form method="POST">
			<div>
				<label>Location Id:</label>
			</div>
			<div>
				<input type="number" name="location_id" value="{{.Main.LocationId}}">
			</div>
			<div>
				<label>API Url:</label>
			</div>
			<div>
				<input type="text" name="api_url" value="{{.API.Url}}">
			</div>
			<div>
				<label>API Id:</label>
			</div>
			<div>
				<input type="text" name="api_id" value="{{.API.Id}}">
			</div>
			<div>
				<label>API Key:</label>
			</div>
			<div>
				<input type="password" name="api_key" value="{{.API.Key}}">
			</div>
			<div>
				<label>XMPP Server:</label>
			</div>
			<div>
				<input type="text" name="xmpp_server" value="{{.XMPP.Server}}">
			</div>
			<div>
				<label>XMPP User:</label>
			</div>
			<div>
				<input type="text" name="xmpp_user" value="{{.XMPP.User}}">
			</div>
			<div>
				<label>XMPP Pass:</label>
			</div>
			<div>
				<input type="password" name="xmpp_pass" value="{{.XMPP.Pass}}">
			</div>
			<div>
				<input type="submit">
			</div>
		</form>
	</body>
</html>
`

var config = `
[Main]
LocationId={{.Main.LocationId}}

[API]
Url={{.API.Url}}
Id={{.API.Id}}
Key={{.API.Key}}

[XMPP]
Server={{.XMPP.Server}}
User={{.XMPP.User}}
Pass={{.XMPP.Pass}}
`

var (
	configTmpl = textTemplate.Must(textTemplate.New("config").Parse(config))
	formTmpl = template.Must(template.New("form").Parse(form))
)

func handle(w http.ResponseWriter, req *http.Request) {
    username, password, _ := req.BasicAuth()

    if username != "admin" || password != "admin" {
    	w.Header().Add("WWW-Authenticate", "Basic")
        http.Error(w, "authorization failed", http.StatusUnauthorized)
        return
    }

    if req.Method == http.MethodGet {
    	Get(w, req)
    } else {
    	Post(w, req)
    }
}

func Get(w http.ResponseWriter, req *http.Request) {
	formTmpl.Execute(w, global.Cfg)
}

func Post(w http.ResponseWriter, req *http.Request) {
	if err := post(w, req); err != nil {
		log.Printf("post: %v", err)
		http.Error(w, "Wrong input values", http.StatusBadRequest)
		return
	}
	io.WriteString(w, "Restarting server...")
	go func() {
		<-time.After(10 * time.Second)
		os.Exit(1)
	}()
}

func post(w http.ResponseWriter, req *http.Request) (err error) {
	locId, err := strconv.ParseInt(req.FormValue("location_id"), 10, 64)
	if err != nil {
		return
	}
	global.Cfg.Main.LocationId = locId
	global.Cfg.API.Url = req.FormValue("api_url")
	global.Cfg.API.Id = req.FormValue("api_id")
	global.Cfg.API.Key = req.FormValue("api_key")
	global.Cfg.XMPP.Server = req.FormValue("xmpp_server")
	global.Cfg.XMPP.User = req.FormValue("xmpp_user")
	global.Cfg.XMPP.Pass = req.FormValue("xmpp_pass")
	f, err := os.Create("conf/gateway.conf")
	if err != nil {
		return
	}
	defer f.Close()
	return configTmpl.Execute(f, global.Cfg)
}



func NewHttp() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
