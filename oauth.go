package robinhood

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// DefaultClientID is used by the website.
const DefaultClientID = "c82SH0WZOsabOXGP2sxqcj34FxkvfnWRZBKlBjFS"

// OAuth implements oauth2 using the robinhood implementation
type OAuth struct {
	Endpoint, ClientID, Username, Password, MFA string
	DeviceID                                    string
}

// ErrMFARequired indicates the MFA was required but not provided.
var ErrMFARequired = fmt.Errorf("Two Factor Auth code required and not supplied")

type RhAuth struct {
	DeviceToken string `json:"device_token"`
	ClientID    string `json:"client_id"`
	ExpiresIn   int    `json:"expires_in"`
	GrantType   string `json:"grant_type"`
	Scope       string `json:"scope"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	AlPk        string `json:"al_pk,omitempty"`
	AlToken     string `json:"al_token,omitempty"`
	MfaCode     string `json:"mfa_code"`
}

// Token implements TokenSource
func (p *OAuth) Token() (*oauth2.Token, error) {
	cliID := p.ClientID
	if cliID == "" {
		cliID = DefaultClientID
	}

	authDtr := RhAuth{
		DeviceToken: p.DeviceID,
		ClientID:    cliID,
		ExpiresIn:   int(24 * time.Hour / time.Second),
		GrantType:   "password",
		Scope:       "internal",
		Username:    p.Username,
		Password:    p.Password,
		MfaCode:     "",
	}

	if p.MFA != "" {
		authDtr.MfaCode = p.MFA
	}

	rData, err := json.Marshal(authDtr)
	req, err := http.NewRequest(
		"POST",
		EPLogin,
		bytes.NewReader(rData),
	)

	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "could not post login")
	}
	defer res.Body.Close()

	var o struct {
		oauth2.Token
		ExpiresIn   int    `json:"expires_in"`
		MFARequired bool   `json:"mfa_required"`
		MFAType     string `json:"mfa_type"`
	}
	err = json.NewDecoder(res.Body).Decode(&o)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode token")
	}

	if o.MFARequired {
		return nil, ErrMFARequired
	}
	o.Token.Expiry = time.Now().Add(time.Duration(o.ExpiresIn) * time.Second)
	return &o.Token, nil
}
