// Package fileupload provides an easy way to upload files to a filehost.
package fileupload

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var hosts = []string{"https://0x0.st", "https://uguu.se/api.php?d=upload-tool"}

// UploadFile takes a file and uploads that file to a file host.
// It returns the url to the uploaded file as a string and any error encountered.
func UploadFile(file *os.File) (string, error) {
	var err error
	var result string

	for _, host := range hosts {
		result, err = UploadToHost(host, file)
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", err
	}
	return result, nil
}

// UploadToHost takes a url and a file as arguments and uploads the file to the provided url with HTTP POST.
// It returns the url to the uploaded file as a string and any error encountered.
func UploadToHost(url string, file *os.File) (string, error) {
	var err error

	values := map[string]io.Reader{
		"file": file,
	}

	var client http.Client
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = writer.CreateFormFile(key, x.Name()); err != nil {
				return "", err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return "", err
		}

	}
	writer.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Submit the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return strings.Replace(bodyString, "\n", "", -1), nil
	}
	return "", fmt.Errorf("bad status: %s", resp.Status)
}
