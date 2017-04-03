package tests

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"runtime"
	//"path/filepath"
	//"runtime"
	"strings"
	"testing"

	"github.com/astaxie/beego"
)

func init() {
	_, file, _, _ := runtime.Caller(1)

	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, "."+string(filepath.Separator))))

	beego.TestBeegoInit(apppath)

	//beego.Run()
}

func makeHttpRequest(method string, path string, body io.Reader) (req *http.Request, write *httptest.ResponseRecorder, err error) {
	if body == nil {
		body = ioutil.NopCloser(strings.NewReader(url.Values{}.Encode()))
	}
	r, e := http.NewRequest(method, path, body)
	//r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	return r, w, e
}

func TestGET(t *testing.T) {

	body := ioutil.NopCloser(strings.NewReader(`{"name": "teste3"}`))

	r, w, _ := makeHttpRequest("GET", "/v1/accounts", body)

	beego.BeeApp.Handlers.ServeHTTP(w, r)

	fmt.Println(w.Body.String())

	//So(w.Body.String(), ShouldContainSubstring, "hi")

}
