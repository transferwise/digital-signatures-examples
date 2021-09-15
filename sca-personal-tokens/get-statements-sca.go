package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	privateKeyPath = "./private.pem" // Change to private key path
	profileID      = ""              // Change to profile ID
	accountID      = ""              // Change to borderless account ID
	baseURL        = "https://api.transferwise.com"
)

var token = os.Getenv("API_TOKEN")

func init() {
	if token == "" {
		panic("no api token, please set with $ export API_TOKEN=xxx")
	}
	if profileID == "" || accountID == "" {
		panic("profile / account ID missing, please add them")
	}
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		panic("private key file not found, please update key path")
	}
}

func main() {
	currency := "GBP"
	periodDays := 14

	statement, err := getStatement(currency, periodDays)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n%s statement received with %d transactions.",
		currency, len(statement.Transactions))
}

// Transaction represents a single statement line, that shows a CREDIT
// or DEBIT movement in the account.
type Transaction struct {
	txnType string    `json:'type,omitempty'`
	date    time.Time `json:'date,omitempty'`
}

// Statement is a period of activity in a currency account.
type Statement struct {
	Transactions []Transaction `json:'transactions,omitempty'`
}

func getStatement(currency string, days int) (Statement, error) {
	statement := &Statement{}
	var oneTimeToken, signature string

	intervalStart := time.Now().AddDate(0, 0, -days).Format(time.RFC3339)
	intervalEnd := time.Now().Format(time.RFC3339)

	params := url.Values{}
	params.Add("currency", currency)
	params.Add("type", "COMPACT")
	params.Add("intervalStart", fmt.Sprintf("%s", intervalStart))
	params.Add("intervalEnd", fmt.Sprintf("%s", intervalEnd))

	response, err := doRequest(oneTimeToken, signature, params)
	if err != nil {
		return Statement{}, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(body, &statement)
	if err != nil {
		return Statement{}, err
	}

	return *statement, err
}

func handleSCA(oneTimeToken string) (string, error) {
	fmt.Println("doing sca challenge")

	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		fmt.Errorf("signer is damaged: %v", err)
	}

	// Hashing the x-2fa-approval token with SHA-256.
	hashed := sha256.Sum256([]byte(oneTimeToken))

	// Signing the token with our private key.
	signedToken, err := rsa.SignPKCS1v15(rand.Reader, privateKey,
		crypto.SHA256, hashed[:],
	)
	if err != nil {
		fmt.Errorf("could not sign token: %v", err)
	}

	// Encoding to string to be included in HTTP headers.
	signature := base64.StdEncoding.EncodeToString(signedToken)

	return signature, nil
}

func loadPrivateKey(filePath string) (*rsa.PrivateKey, error) {
	privateKey := &rsa.PrivateKey{}

	// Read private key .pem file into bytes.
	pemBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &rsa.PrivateKey{}, err
	}

	// Decode the key block from the bytes, ignore rest.
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return &rsa.PrivateKey{}, errors.New("no key found")
	}

	// Detect the key type and parse it into an rsa key object.
	switch block.Type {
	case "RSA PRIVATE KEY":
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return &rsa.PrivateKey{}, err
		}
		return privateKey, nil
	default:
		return &rsa.PrivateKey{}, fmt.Errorf("unsupported key type %q", block.Type)
	}

	return privateKey, nil
}

func doRequest(oneTimeToken, signature string, params url.Values) (*http.Response, error) {
	u := fmt.Sprintf("%s/v3/profiles/%s/borderless-accounts/%s/statement.json",
		baseURL, profileID, accountID)

	url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	url.RawQuery = params.Encode()

	fmt.Printf("\n%s %s\n", http.MethodGet, url.String())

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("User-Agent", "tw-statements-sca")
	request.Header.Set("Content-Type", "application/json")
	if oneTimeToken != "" {
		request.Header.Set("x-2fa-Approval", oneTimeToken)
		request.Header.Set("X-Signature", signature)
		fmt.Printf("\nx-2fa-approval: %s", request.Header.Get("x-2fa-approval"))
		fmt.Printf("\nX-Signature: %s\n", request.Header.Get("X-Signature"))
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	fmt.Println(response.StatusCode)

	switch response.StatusCode {
	case http.StatusOK:
		return response, nil
	case http.StatusForbidden:
		if response.Header.Get("x-2fa-approval") != "" {
			oneTimeToken = response.Header.Get("x-2fa-approval")
			signature, err = handleSCA(oneTimeToken)
			if err != nil {
				return nil, err
			}
			return doRequest(oneTimeToken, signature, params)
		}
	default:
		return nil, fmt.Errorf("\nsomething failed, http status: %d", response.StatusCode)
	}

	return response, nil
}
