package kit

import (
	"fmt"
	"os"

	"github.com/imroc/req/v3"
)

type APIError struct {
	Message          string `json:"message"`
	DocumentationUrl string `json:"documentation_url"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error: %s (refer to %s)", e.Message, e.DocumentationUrl)
}
func AfterResponse(client *req.Client, resp *req.Response) error {
	if resp.Err != nil { // resp.Err represents the underlying error, e.g. network error, or unmarshal error (SetResult or SetError was invoked before).
		if dump := resp.Dump(); dump != "" { // Append dump content to original underlying error to help troubleshoot if request has been sent.
			resp.Err = fmt.Errorf("%s\nraw content:\n%s", resp.Err.Error(), resp.Dump())
		}
		return nil // Skip the following logic if there is an underlying error.
	}
	// Return a human-readable error if server api returned an error message.
	if err, ok := resp.ErrorResult().(*APIError); ok {
		resp.Err = err
		return nil
	}
	// Corner case: neither an error response nor a success response (e.g. status code < 200),
	// dump content to help troubleshoot.
	if !resp.IsSuccessState() {
		resp.Err = fmt.Errorf("bad response, raw content:\n%s", resp.Dump())
		return nil
	}
	return nil
}

func NewClient() *req.Client {
	c := req.C()
	if os.Getenv("DEBUG") != "" {
		c = c.DevMode()
	}
	c.OnAfterResponse(AfterResponse)
	return c
}

func ClientPost[IN, OUT any](c *req.Client, path string, in IN, out *OUT) error {
	resp := c.Post(path).SetBody(in).Do()
	if resp.IsErrorState() {
		return resp.Err
	}
	err := resp.Into(out)
	if err != nil {
		return resp.Err
	}
	return nil
}
