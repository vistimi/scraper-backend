package host

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type Request struct {
	Host   string
	Args   map[string]string
	Header map[string][]string
}

type nopCloser struct {
	io.Reader
}

// Wrapper for Close of io reader
func (nopCloser) Close() error { return nil }

// Generate the URL based on the request parameters
func (request *Request) URL() string {
	args := request.Args
	s := request.Host + EncodeQuery(args)
	return s
}

// Send http request and read the body of the response
func (request *Request) ExecuteGET() (response []byte, ret error) {

	url := request.URL()
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = request.Header
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	return body, nil
}

// Generate a string buffer based on parameters
func EncodeQuery(args map[string]string) string {
	i := 0
	s := bytes.NewBuffer(nil)
	for k, v := range args {
		if i != 0 {
			s.WriteString("&")
		}
		i++
		s.WriteString(k + "=" + url.QueryEscape(v))
	}
	return s.String()
}

// GetFile returns the buffer of the downloaded image
func GetFile(URL string) ([]byte, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("Received non 200 response code")
	}

	buffer, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll has failed: %v", err)
	}
	return buffer, nil
}

// Download a file from an URL body response
func DownloadFile(URL string, fileName string) error {
	//Get the response bytes from the url
	buffer, err := GetFile(URL)

	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, bytes.NewReader(buffer))
	if err != nil {
		return err
	}
	return nil
}
