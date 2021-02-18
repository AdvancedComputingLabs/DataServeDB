package runtime

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AdvancedComputingLabs/httpserv"
)

const httpPort = 8080
const httpsPort = 10443

const REDIRECT_TO_HTTPS = false

var httpsHost string
var servers []*httpserv.HttpServer

func StartHttpServer() {
	//NOTE: don't need starting http/https message as it is handled in httpserv package.

	//NOTE: Redirect the incoming HTTP request. Note that "127.0.0.1:8081" will only work if you are accessing the server from your local machine.

	//TODO: change local host
	httpsHost = fmt.Sprintf("https://localhost:%d", httpsPort)

	//configDir := paths.GetConfigDirPath()
	//TODO: load names from config file.
	//servCert := paths.Combine(configDir, "server.crt")
	//servKey := paths.Combine(configDir, "server.key")

	//TODO: error handling?
	servers = append(servers, httpserv.NewHttpServ(httpPort, "/", httpServReqHandler))
	//servers = append(servers, httpserv.NewHttpServTLS(httpsPort, "/", httpsServReqHandler, servCert, servKey))

	//timed break for printing correctly.
	time.Sleep(1)

	//TODO:
	//	NEXT: shutdown message handling.

	//TODO: move it http server package.
	//log.Printf("Closing http and https servers ...")
}

func httpServReqHandler(w http.ResponseWriter, r *http.Request) {
	if REDIRECT_TO_HTTPS {
		redirectToHttps(w, r)
		return
	}

	commonHttpServReqHandler(w, r)
}

func httpsServReqHandler(w http.ResponseWriter, r *http.Request) {
	commonHttpServReqHandler(w, r)
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, httpsHost+r.RequestURI, http.StatusMovedPermanently)
}
