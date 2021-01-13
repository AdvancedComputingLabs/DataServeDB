// Copyright (c) 2018 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

package httpserv

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HttpServer struct {

	/* Private fields protected against outside tempering. */
	certFile              string
	errorOnInitialization error
	finished              chan struct{}
	handler               func(http.ResponseWriter, *http.Request)
	initialized           bool
	startedForTLS         bool
	keyFile               string
	port                  int
	rootDir               string
	signalStopNotifier    chan struct{}
}

func createHttpServer(port int, root_dir string, request_handler_fptr func(http.ResponseWriter, *http.Request)) *HttpServer {
	srv := HttpServer{
		port:               port,
		handler:            request_handler_fptr,
		finished:           make(chan struct{}),
		rootDir:            strings.TrimSpace(root_dir),
		signalStopNotifier: make(chan struct{}),
	}

	if len(srv.rootDir) == 0 {
		srv.rootDir = "/"
	}

	return &srv
}

func (t *HttpServer) GetPort() int {
	return t.port
}

//This is for documentations only.
//All functions ending with async must be called as go routine.
//All async functions must have failed callback handler or channel as they cannot return error.
func (t *HttpServer) initHttpsServerAsync() {
	//TODO: add callback mechanism to handle errors/failures.
	//TODO: Test tls for cert related failures. Does it start as http or fails?
	//NOTE: log.Printf shouldn't be here, but this is experimental code so it is ok.

	if t.startedForTLS {
		log.Printf("Starting https server on port %d ...\n", t.port)
	} else {
		log.Printf("Starting http server on port %d ...\n", t.port)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(t.rootDir, t.handler)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", t.port),
		Handler: mux}

	idleConnsClosed := make(chan struct{})

	go func() {
		<-t.signalStopNotifier

		if t.startedForTLS {
			log.Printf("Shutting https server on port %d ...\n", t.port)
		} else {
			log.Printf("Shutting http server on port %d ...\n", t.port)
		}

		//Gracefully shutdown http server.
		if err := srv.Shutdown(context.Background()); err != nil {
			//TODO: error notification for parent process. Could add errorOnShutdown to the HttpServer.
			// Error from closing listeners, or context timeout:
			log.Printf("Error during http(s) server shutdown: %v", err)
		}

		close(t.finished)
		close(idleConnsClosed)

	}()

	var err error

	if t.startedForTLS {
		err = srv.ListenAndServeTLS(t.certFile, t.keyFile)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		t.errorOnInitialization = err

		// Error starting or closing listener:
		if t.startedForTLS {
			log.Printf("HTTPS server ListenAndServe: %v", err)
		} else {
			log.Printf("HTTP server ListenAndServe: %v", err)
		}
	}

	<-idleConnsClosed
}

func (t HttpServer) IsInitialized() (bool, error) {
	return t.initialized, t.errorOnInitialization
}

func NewHttpServ(port int, root_dir string, request_handler_fptr func(http.ResponseWriter, *http.Request)) *HttpServer {
	srv := createHttpServer(port, root_dir, request_handler_fptr)
	go srv.initHttpsServerAsync()
	return srv
}

func NewHttpServTLS(port int, root_dir string, request_handler_fptr func(http.ResponseWriter, *http.Request), certFile string, keyFile string) *HttpServer {

	srv := createHttpServer(port, root_dir, request_handler_fptr)

	srv.startedForTLS = true
	srv.certFile = certFile
	srv.keyFile = keyFile

	go srv.initHttpsServerAsync()

	return srv
}

func (t *HttpServer) signalStop() {
	close(t.signalStopNotifier)
}

func (t *HttpServer) Shutdown() {
	t.signalStop()
	t.waitToFinish()
}

//Function doesn't guarantee server is running for tls but shows if server was started with NewHttpServTLS;
//it fallbacks to the default behavior of go's http tls procedure.
func (t HttpServer) StartedForTLS() bool {
	return t.startedForTLS
}

func (t *HttpServer) waitToFinish() {
	<-t.finished
}
