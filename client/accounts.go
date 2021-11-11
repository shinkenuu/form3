package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

type Form3Client struct {
	// Api address ([scheme://]host[:port])
	apiAddress string

	// Api path to account resource
	accountPath string

	// HTTP client to make requests with
	httpClient *http.Client
}

func New() *Form3Client {
	apiAddress := os.Getenv("API_ADDR")
	if apiAddress == "" {
		apiAddress = `http://127.0.0.1:8080`
	}

	client := Form3Client{
		apiAddress:  apiAddress,
		accountPath: "v1/organisation/accounts",
		httpClient:  &http.Client{},
	}

	return &client
}

func (f *Form3Client) FetchAccount(ID string) (*AccountData, error) {
	resourceUrl, err := f.accountResourceUrl(ID)
	if err != nil {
		return nil, err
	}

	responseBody, err := f.doRequest(resourceUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}

	var fetchedAccount Account
	err = json.Unmarshal(responseBody, &fetchedAccount)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return fetchedAccount.AccountData, nil
}

func (f *Form3Client) CreateAccount(accountData *AccountData) (*AccountData, error) {
	resourceUrl, err := f.accountResourceUrl("")
	if err != nil {
		return nil, err
	}

	account := Account{AccountData: accountData}
	jsonBody, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	responseBody, err := f.doRequest(resourceUrl, "POST", bytes.NewBuffer(jsonBody), nil)
	if err != nil {
		return nil, err
	}

	var createdAccount Account
	err = json.Unmarshal(responseBody, &createdAccount)
	if err != nil {
		return nil, err
	}

	return createdAccount.AccountData, err
}

func (f *Form3Client) DeleteAccount(ID string, version int64) error {
	resourceUrl, err := f.accountResourceUrl(ID)

	query := map[string]string{"version": strconv.FormatInt(version, 10)}

	if err != nil {
		return err
	}

	_, err = f.doRequest(resourceUrl, "DELETE", nil, query)
	return err
}

func (f *Form3Client) doRequest(url string, method string, body io.Reader, query map[string]string) ([]byte, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	setHeaders(request)
	setQuery(request, query)

	response, err := f.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	log.Println("Response status: " + response.Status)
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, errors.New(string(responseBody))
	}

	log.Println("Response body: ", string(responseBody))
	return responseBody, nil
}

// Set API required headers at `request`
func setHeaders(request *http.Request) {
	request.Header.Set("Accept", "application/vnd.api+json")
	request.Header.Set("Host", "api.form3.tech")
	request.Header.Set("Date", time.Now().UTC().String())
	request.Header.Set("Content-Type", "application/vnd.api+json")
}

// Set `request`'s query with `query`'s keys and values
func setQuery(request *http.Request, query map[string]string) {
	requestQuery := request.URL.Query()

	for key, value := range query {
		requestQuery.Add(key, value)
	}

	request.URL.RawQuery = requestQuery.Encode()
}

// Format an URL for a Account resource
func (f *Form3Client) accountResourceUrl(resourceID string) (string, error) {
	parsedUrl, err := url.Parse(f.apiAddress)
	if err != nil {
		return "", err
	}

	parsedUrl.Path = path.Join(parsedUrl.Path, f.accountPath, resourceID)
	return parsedUrl.String(), nil
}
