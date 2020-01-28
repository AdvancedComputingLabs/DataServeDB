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

package unstable_api

import (
	"fmt"
	"net/http"
)

func aclabsHttpServReqHandler(w http.ResponseWriter, r *http.Request) {
	if REDIRECT_TO_HTTPS {
		redirectToHttps(w, r)
		return
	}

	commonHttpServReqHandler(w, r)
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	//NOTE: Redirect the incoming HTTP request. Note that "127.0.0.1:8081" will only work if you are accessing the server from your local machine.
	http.Redirect(w, r, fmt.Sprintf("https://localhost:%d", httpsPort)+r.RequestURI, http.StatusMovedPermanently)
}
