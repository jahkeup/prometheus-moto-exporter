package gather

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jahkeup/prometheus-moto-exporter/pkg/hnap"
)

const hSOAPAction = "SOAPAction"
const hHNAPAuth = "HNAP_AUTH"

type Gatherer struct {
	username string
	password string

	endpoint *url.URL

	mu         *sync.RWMutex
	privateKey []byte
	client     *http.Client
}

func New(endpoint *url.URL, username, password string) (*Gatherer, error) {
	return &Gatherer{
		username: username,
		password: password,
		endpoint: endpoint,

		mu: &sync.RWMutex{},

		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Jar: func() http.CookieJar { j, _ := cookiejar.New(nil); return j }(),
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: time.Second * 45,
		},
	}, nil
}

func (g *Gatherer) Login() error {
	const (
		loginAction = "Login"
		loginURI    = "http://purenetworks.com/HNAP1/Login"
	)

	log := logrus.WithField("action", "login")

	// 1. Request challenge, uid, and public key from endpoint. We have to use a
	// valid username to be given a login challenge.

	challenge := map[string]interface{}{
		// Wrap the message in the HNAP action name.
		"Login": map[string]string{
			"Action":   "request",
			"Username": g.username,
		},
	}
	data, err := json.Marshal(challenge)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, g.endpoint.String(), bytes.NewReader(data))
	req.Header.Add(hSOAPAction, loginURI)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	log.Debug("requesting challenge")
	resp, err := g.client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("unable to request challenge")
		return err
	}
	log.Debug("accepting challenge")

	hnapResponse := struct {
		LoginResponse struct {
			Challenge string
			PublicKey string
			// Should be held onto by the shared session.
			Cookie string
		}
	}{}

	err = json.NewDecoder(resp.Body).Decode(&hnapResponse)
	if err != nil {
		return err
	}
	log.WithFields(logrus.Fields{
		"challenge": hnapResponse.LoginResponse.Challenge,
		"uid":       hnapResponse.LoginResponse.Cookie,
	}).Trace("computing response")

	// 2. Compute challenge response by making its "private key". We'll use it
	// to submit a login challenge response to complete the login-flow.

	privateKey, err := digest(hnapResponse.LoginResponse.Challenge, []byte(hnapResponse.LoginResponse.PublicKey+g.password))
	if err != nil {
		return err
	}

	passKey, err := digest(hnapResponse.LoginResponse.Challenge, privateKey)
	if err != nil {
		return err
	}

	uidCookie := &http.Cookie{
		Name:  "uid",
		Value: string(hnapResponse.LoginResponse.Cookie),
	}
	pkCookie := &http.Cookie{
		Name:  "PrivateKey",
		Value: string(privateKey),
	}

	// 3. Submit response to challenge to complete the login.

	login := map[string]interface{}{
		"Login": map[string]string{
			"Action":        "login",
			"Username":      g.username,
			"LoginPassword": string(passKey),
		},
	}

	data, err = json.Marshal(login)
	if err != nil {
		return err
	}

	req, err = g.requestWithKey(loginAction, loginURI, bytes.NewReader(data), privateKey)
	if err != nil {
		return err
	}

	req.AddCookie(uidCookie)
	req.AddCookie(pkCookie)

	log.Debug("submitting response")
	resp, err = g.client.Do(req)
	if err != nil {
		log.WithError(err).Error("unable to login")
		return err
	}
	resp.Body.Close()

	log.WithFields(logrus.Fields{
		"action":      loginURI,
		"action.call": "login",
		"status":      resp.StatusCode,
	}).Debug("response sent")

	if resp.StatusCode != http.StatusOK {
		return errors.New("challenge response rejected")
	}

	// Update client to use our new session

	log.Trace("updating gatherer HTTP client")
	// Acquire lock to modify the underlying client data.
	g.mu.Lock()
	{
		// Record the Private Key that's for this login.
		g.privateKey = privateKey
		g.client.Jar.SetCookies(g.endpoint, []*http.Cookie{uidCookie, pkCookie})
	}
	g.mu.Unlock()
	log.Trace("gatherer configured with new login session")

	return nil
}

func (g *Gatherer) Gather() (*Collection, error) {
	const actionName = hnap.GetMultipleHNAPs
	const actionURI = "http://purenetworks.com/HNAP1/" + actionName

	log := logrus.WithField("action", actionURI)

	data, err := json.Marshal(hnap.GetMultipleRequestData(
		hnap.GetHomeAddress,
		hnap.GetHomeConnection,
		hnap.GetMotoLagStatus,
		hnap.GetMotoStatusConnectionInfo,
		hnap.GetMotoStatusDownstreamChannelInfo,
		hnap.GetMotoStatusLog,
		hnap.GetMotoStatusSoftware,
		hnap.GetMotoStatusStartupSequence,
		hnap.GetMotoStatusUpstreamChannelInfo,
	))
	if err != nil {
		return nil, err
	}

	g.mu.RLock()
	unlock := unlockGuarded(g.mu.RLocker())
	defer unlock()

	req, err := g.request(actionName, actionURI, bytes.NewReader(data))
	if err != nil {
		log.Error("unable to prepare request")
		return nil, err
	}
	resp, err := g.client.Do(req)
	if err != nil {
		log.WithError(err).Error("unable to complete request")
		return nil, err
	}
	unlock()

	defer resp.Body.Close()

	var response hnap.GetMultipleHNAPsResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	for k, v := range response.HNAP {
		// Raw JSON string
		logrus.WithField("name", k).Trace("%s", v)
	}

	data, err = response.GetJSON(hnap.GetMotoStatusDownstreamChannelInfo)
	var downstream hnap.DownstreamChannelResponse
	err = json.Unmarshal(data, &downstream)
	if err != nil {
		return nil, err
	}

	data, err = response.GetJSON(hnap.GetMotoStatusUpstreamChannelInfo)
	var upstream hnap.UpstreamChannelResponse
	err = json.Unmarshal(data, &upstream)
	if err != nil {
		return nil, err
	}

	return &Collection{
		Upstream: upstream.Channels,
		Downstream: downstream.Channels,
	}, nil
}

func (g *Gatherer) request(actionName, actionURI string, data io.Reader) (*http.Request, error) {
	return g.requestWithKey(actionName, actionURI, data, g.privateKey)
}

func (g *Gatherer) requestWithKey(actionName, actionURI string, data io.Reader, key []byte) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, g.endpoint.String(), data)
	if err != nil {
		return nil, err
	}

	hnapAuth, ts, err := digestAuth(actionURI, key)
	if err != nil {
		return nil, err
	}

	req.Header.Add(hSOAPAction, fmt.Sprintf(`"%s"`, actionURI))
	req.Header.Add(hHNAPAuth, fmt.Sprintf("%s %d", string(hnapAuth), ts))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	return req, nil
}

// Single use Unlock()er.
func unlockGuarded(lock sync.Locker) func() {
	var singleUse sync.Once
	return func() { singleUse.Do(lock.Unlock) }
}

// digestAuth prepares an authentication digest for calling a given SOAPAction.
func digestAuth(actionURI string, key []byte) ([]byte, int64, error) {
	ts := time.Now().Unix()
	data, err := digest(fmt.Sprintf(`%d"%s"`, ts, actionURI), key)
	return data, ts, err
}

// digest prepares an authentication digest for use with HNAP.
func digest(msg string, key []byte) ([]byte, error) {
	mac := hmac.New(md5.New, key)
	_, err := fmt.Fprintf(mac, msg)
	if err != nil {
		return nil, err
	}
	digestData := mac.Sum(nil)

	digestHex := make([]byte, hex.EncodedLen(len(digestData)))
	hex.Encode(digestHex, digestData)
	return bytes.ToUpper(digestHex), nil
}
