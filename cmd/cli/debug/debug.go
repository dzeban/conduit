package debug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"
	"github.com/tidwall/pretty"

	"github.com/dzeban/conduit/cmd/cli/state"
)

func MakeRequestWithDump(method, url string, data interface{}) (*http.Response, []byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to marshal json")
	}

	bb := bytes.NewBuffer(b)

	req, err := http.NewRequest(method, url, bb)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create http request")
	}
	req.Header.Add("Content-Type", "application/json")

	reqDump, err := httputil.DumpRequest(req, false)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to dump http request")
	}

	fmt.Println("================================")
	fmt.Printf("%s", reqDump)
	fmt.Printf("%s\n", pretty.Color(pretty.Pretty(b), nil))
	fmt.Println("--------------------------------")

	resp, err := state.Client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to make http request")
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to read http response body")
	}

	respDump, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to dump http response")
	}

	fmt.Printf("%s", respDump)
	fmt.Printf("%s\n", pretty.Color(pretty.Pretty(body), nil))
	fmt.Println("================================")

	return resp, body, nil
}

func MakeAuthorizedRequestWithDump(method, url string, data interface{}) (*http.Response, []byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to marshal json")
	}

	bb := bytes.NewBuffer(b)

	req, err := http.NewRequest(method, url, bb)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create http request")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+state.CurrentToken)

	reqDump, err := httputil.DumpRequest(req, false)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to dump http request")
	}

	fmt.Println("================================")
	fmt.Printf("%s", reqDump)
	fmt.Printf("%s\n", pretty.Color(pretty.Pretty(b), nil))
	fmt.Println("--------------------------------")

	resp, err := state.Client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to make http request")
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to read http response body")
	}

	respDump, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to dump http response")
	}

	fmt.Printf("%s", respDump)
	fmt.Printf("%s\n", pretty.Color(pretty.Pretty(body), nil))
	fmt.Println("================================")

	return resp, body, nil
}
