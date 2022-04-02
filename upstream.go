package dirtyhttp

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpClientInterface interface {
	Do(*http.Request) ([]byte, error)
}

type httpClient struct{}

func (c *httpClient) Do(r *http.Request) ([]byte, error) {
	return httpDo(r)
}

func httpDo(r *http.Request) ([]byte, error) {
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		return []byte{}, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func (c *httpClient) DoForStatus(r *http.Request, status int) ([]byte, error) {
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

    if res.StatusCode != status {
        return body, errors.New(fmt.Sprintf("Status code was %d, instead of the expected %d", res.StatusCode, status))
    } else {
        return body, nil
    }

}

type upstream struct {
    Db *sql.DB
    Http HttpClientInterface
}

func newUpstream() *upstream {
    return &upstream{Http: &httpClient{}}
}

func (u *upstream) SetDatabase(database *sql.DB) {
	u.Db = database
}

func (u *upstream) SetHttpClient(client HttpClientInterface) {
    u.Http = client
}


