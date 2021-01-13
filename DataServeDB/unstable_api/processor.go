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

// main.go project main.go

package unstable_api

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	/*aclabs packages seperated*/
	//aclabs_httpserv "github.com/AdvancedComputingLabs/httpserv"
	// aclabs_httpserv "./httpserv"
)

const httpPort = 8080
const httpsPort = 10443

const REDIRECT_TO_HTTPS = true

var servers []*HttpServer

// var cache redis.Conn
var cache Store

func cliProcessor() {
	reader := bufio.NewReader(os.Stdin)
	keep_running := true

	for keep_running {

		fmt.Println("Enter Command: ")
		cmd_text, _ := reader.ReadString('\n')
		cmd_text = strings.Trim(cmd_text, "\r\n")
		cmd_text_toks := strings.Split(cmd_text, " ")

		if len(cmd_text_toks) > 0 {
			switch strings.ToUpper(cmd_text_toks[0]) {

			case "EXIT":
				for _, s := range servers {
					s.SignalStop()
					s.WaitToFinish()
				}
				keep_running = false
			}
		}

	}
}

// func initCache() {
// 	// Initialize the redis connection to a redis instance running on your local machine
// 	conn, err := redis.DialURL("redis://localhost")
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Assign the connection to the package level `cache` variable
// 	// cache = conn
// }

func Process() {
	log.Println("Starting DataServe...")
	// initCache()
	cache = Store{}
	cache.Init()

	// servers = append(servers, aclabs_httpserv.NewHttpServ(httpPort, "/signin", Signin))
	// servers = append(servers, aclabs_httpserv.NewHttpServTLS(httpsPort, "/signin", Signin, "certfiles/server.crt", "certfiles/server.key"))
	servers = append(servers, NewHttpServ(httpPort, "/", aclabsHttpServReqHandler))
	servers = append(servers, NewHttpServTLS(httpsPort, "/", aclabsHttpsServReqHandler, "unstable_api/certfiles/server.crt", "unstable_api/certfiles/server.key"))

	//timed break for printing correctly.
	time.Sleep(1)

	//problem in test env
	//cliProcessor()

	log.Println("Closing DataServe...")
}
