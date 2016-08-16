package endpoint

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func formUpload(ip net.IP, path, fileName, fileContent string,
	noReply bool) (string, error) {
	u := url.URL{
		Scheme: "http",
		Host:   ip.String(),
		Path:   path,
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("1", fileName)
	if err != nil {
		return "", err
	}
	if _, err = fw.Write([]byte(fileContent)); err != nil {
		return "", err
	}

	w.Close()

	req, err := http.NewRequest("POST", u.String(), &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		if strings.HasSuffix(err.Error(), "EOF") && noReply {
			return "", nil
		}

		return "", err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
