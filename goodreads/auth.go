package goodreads

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
)

const (
	cookieDomain  = ".goodreads.com"
	authDomain    = "goodreads.com"
	authBaseUrl   = "https://api.amazon.com/auth"
	authUserAgent = "AmazonWebView/GoodreadsForIOS App/4.9.0/iOS/16.6.1/iPhone"
)

type frcCookieBody struct {
	ApplicationVersion         string
	DeviceLanguage             string
	DeviceFingerprintTimestamp string
	DeviceOSVersion            string
	DeviceName                 string
	ScreenHeightPixels         string
	ThirdPartyDeviceId         string
	TimeZone                   string
	ApplicationName            string
	ScreenWidthPixels          string
	DeviceJailbroken           bool
}

type RegisterRequest struct {
	RequestedExtensions []string `json:"requested_extensions"`

	Cookies          Cookies          `json:"cookies"`
	RegistrationData RegistrationData `json:"registration_data"`
	AuthData         AuthData         `json:"auth_data"`

	UserContextMap struct {
		Frc string `json:"frc"`
	} `json:"user_context_map"`

	RequestedTokenType []string `json:"requested_token_type"`
}

type Cookies struct {
	WebsiteCookies []string `json:"website_cookies"`
	Domain         string   `json:"domain"`
}

type RegistrationData struct {
	Domain          string `json:"domain"`
	AppVersion      string `json:"app_version"`
	DeviceType      string `json:"device_type"`
	OsVersion       string `json:"os_version"`
	DeviceSerial    string `json:"device_serial"`
	DeviceModel     string `json:"device_model"`
	AppName         string `json:"app_name"`
	SoftwareVersion string `json:"software_version"`
}

type AuthData struct {
	UserIdPassword struct {
		UserId   string `json:"user_id"`
		Password string `json:"password"`
	} `json:"user_id_password"`
}

type RegisterResponse struct {
	Response struct {
		Success struct {
			Tokens struct {
				Bearer struct {
					AccessToken  string `json:"access_token"`
					RefreshToken string `json:"refresh_token"`
					ExpiresIn    string `json:"expires_in"`
				} `json:"bearer"`
			} `json:"tokens"`
		} `json:"success"`
	} `json:"response"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *Client) Register(ctx context.Context, email, password string) (*RegisterResponse, error) {
	req, err := c.registerRequest(ctx, email, password)
	if err != nil {
		return nil, err
	}

	resp := RegisterResponse{}
	err = c.apiClient.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	c.accessToken = resp.Response.Success.Tokens.Bearer.AccessToken
	c.refreshToken = resp.Response.Success.Tokens.Bearer.RefreshToken

	return &resp, nil
}

func (c *Client) Token(ctx context.Context) (*TokenResponse, error) {
	req, err := c.tokenRequest(ctx)
	if err != nil {
		return nil, err
	}

	resp := TokenResponse{}
	err = c.apiClient.Do(req, &resp)
	return &resp, err
}

func (c *Client) registerRequest(ctx context.Context, email, password string) (*http.Request, error) {
	deviceSerial, err := generateDeviceSerial()
	if err != nil {
		return nil, err
	}

	frcCookie, err := generateFrcCookie(deviceSerial)
	if err != nil {
		return nil, err
	}

	body := RegisterRequest{
		RequestedExtensions: []string{"device_info", "customer_info"},
		Cookies: Cookies{
			Domain: cookieDomain,
		},
		RegistrationData: RegistrationData{
			Domain:          "Device",
			AppVersion:      "4.9.0",
			DeviceType:      "A3NWHXTQ4EBCZS",
			OsVersion:       "16.6.1",
			DeviceSerial:    deviceSerial,
			DeviceModel:     "iPhone",
			AppName:         "GoodreadsForIOS App",
			SoftwareVersion: "1",
		},
		RequestedTokenType: []string{"bearer", "mac_dms", "website_cookies"},
	}

	body.AuthData.UserIdPassword.UserId = email
	body.AuthData.UserIdPassword.Password = password
	body.UserContextMap.Frc = frcCookie

	return c.authRequest(ctx, "POST", authBaseUrl+"/register", body)
}

func (c *Client) tokenRequest(ctx context.Context) (*http.Request, error) {
	val := url.Values{}
	val.Set("app_name", "GoodreadsForIOS App")
	val.Set("app_version", "4.9.0")
	val.Set("di.sdk.version", "6.12.1")
	val.Set("source_token", c.refreshToken)
	val.Set("package_name", "com.goodreads.Goodreads")
	val.Set("di.hw.version", "iPhone")
	val.Set("platform", "iOS")
	val.Set("requested_token_type", "access_token")
	val.Set("source_token_type", "refresh_token")
	val.Set("di.os.name", "iOS")
	val.Set("di.os.version", "16.6.1")
	val.Set("current_version", "6.12.1")
	val.Set("previous_version", "6.12.1")

	return c.authRequest(ctx, "POST", authBaseUrl+"/token", val)
}

func (c *Client) authRequest(ctx context.Context, method, url string, body any) (*http.Request, error) {
	req, err := c.apiClient.Request(ctx, method, url, body)
	req.Header.Set("User-Agent", authUserAgent)
	req.Header.Set("x-amzn-identity-auth-domain", authDomain)

	return req, err
}

func generateDeviceSerial() (string, error) {
	rnd := make([]byte, 16)
	_, err := rand.Read(rnd)
	return fmt.Sprintf("%X", rnd), err
}

func generateUuid() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return strings.ToUpper(id.String()), nil
}

func generateFrcCookie(deviceSerial string) (string, error) {
	now := time.Now()
	ts := fmt.Sprintf("%d", now.UnixMilli())
	tz := now.Format("-07:00")

	deviceUuid, err := generateUuid()
	if err != nil {
		return "", err
	}

	body := frcCookieBody{
		ApplicationVersion:         "4.9.0",
		DeviceLanguage:             "en-US",
		DeviceFingerprintTimestamp: ts,
		DeviceOSVersion:            "iOS/16.6.1",
		DeviceName:                 "iPhone",
		ScreenHeightPixels:         "844",
		ThirdPartyDeviceId:         deviceUuid,
		TimeZone:                   tz,
		ApplicationName:            "Goodreads",
		ScreenWidthPixels:          "390",
		DeviceJailbroken:           false,
	}

	marshaled, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	compressed, err := compress(marshaled)
	if err != nil {
		return "", err
	}

	encrypted, err := frcEncrypt(deviceSerial, compressed)
	if err != nil {
		return "", err
	}

	mac := frcHmac(deviceSerial, encrypted)

	data := []byte{0}
	data = append(data, mac...)
	data = append(data, encrypted...)

	encoded := base64.StdEncoding.EncodeToString(data)

	return encoded, nil
}

func compress(p []byte) ([]byte, error) {
	buf := &bytes.Buffer{}

	w, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(p)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	return buf.Bytes(), err
}

func frcEncrypt(deviceSerial string, src []byte) ([]byte, error) {
	key := pbkdf2.Key(
		[]byte(deviceSerial),
		[]byte("AES/CBC/PKCS7Padding"),
		1000, 16, sha1.New,
	)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padded := pkcs7Pad(src, aes.BlockSize)

	dst := make([]byte, aes.BlockSize+len(padded))
	iv := dst[:aes.BlockSize]

	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}

	enc := cipher.NewCBCEncrypter(block, iv)
	enc.CryptBlocks(dst[aes.BlockSize:], padded)

	return dst, nil
}

func frcHmac(deviceSerial string, data []byte) []byte {
	key := pbkdf2.Key(
		[]byte(deviceSerial),
		[]byte("HmacSHA256"),
		1000, 16, sha1.New,
	)

	h := hmac.New(sha256.New, key)
	h.Write(data)
	sum := h.Sum(nil)
	return sum[:8]
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padSize := blockSize - (len(data) % blockSize)
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(data, pad...)
}
