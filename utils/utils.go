package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
)

// NewClient returns a new client connection.
func NewClient(socketPath string) (*httputil.ClientConn, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to dial socket %q: %v", socketPath, err)
	}
	return httputil.NewClientConn(conn, nil), nil
}

// Do queries the URL and returns the body in bytes.
func Do(client *httputil.ClientConn, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request object: %v", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query endpoint %q: %v", url, err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}
	return b, nil
}

// WriteSnapshot writes raw data to the filesystem.
func WriteSnapshot(dest, name string, data []byte) (int, error) {
	path := filepath.Join(dest, name)
	f, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open %q for writing: %v", path, err)
	}
	defer f.Close()
	n, err := f.Write(data)
	if err != nil {
		return 0, fmt.Errorf("failed to write %q: %v", path, err)
	}
	f.Sync()
	return n, nil
}
