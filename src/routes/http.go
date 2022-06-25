package routes

import (
	"io"
	"net/http"
	"errors"
	"os"
	"bytes"
	"net/url"
	"io/ioutil"
)

type Request struct {
	Host string
	Args   map[string]string
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
func (request *Request) Execute() (response []byte, ret error) {
	s := request.URL()

	res, err := http.Get(s)
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

// Download a file from an URL body response
func DownloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}