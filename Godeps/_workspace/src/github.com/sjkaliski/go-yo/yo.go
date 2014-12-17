// The yo package provides an api client for the Yo app api,
// with methods to send Yo's to users.
package yo

import (
	"errors"
	"net/http"
	"net/url"
)

// Yo API endpoint.
var YO_API = "http://api.justyo.co"

// Yo API Client.
type Client struct {
	Token string
}

// Creates a new Client.
func NewClient(token string) *Client {
	return &Client{
		Token: token,
	}
}

// Sends a "Yo" to all users who subscribe to the active
// account. Expects a 201 response.
func (c *Client) YoAll() error {
	data := url.Values{}
	data.Set("api_token", c.Token)
	res, err := http.PostForm(YO_API+"/yoall/", data)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return errors.New("Received response with non 201 status code.")
	}

	return nil
}

// YoAllLink sends a "Yo" to all subscribed users of the API account
// with the link specified.
func (c *Client) YoAllLink(link string) error {
	data := url.Values{}
	data.Set("api_token", c.Token)
	data.Set("link", link)

	res, err := http.PostForm(YO_API+"/yoall/", data)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return errors.New("Received response with non 201 status code.")
	}

	return nil
}

// Sends a "Yo" to the specified user (who must subscribe)
// to the active account. Expects a 201 response.
func (c *Client) YoUser(username string) error {
	data := url.Values{}
	data.Set("api_token", c.Token)
	data.Set("username", username)
	res, err := http.PostForm(YO_API+"/yo/", data)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return errors.New("Received response with non 201 status code.")
	}

	return nil
}

// YoUserLink sends a "Yo" to the specified user (who must subscribe) with a link
// to the active account. Expects a 201 response.
func (c *Client) YoUserLink(username, link string) error {
	data := url.Values{}
	data.Set("api_token", c.Token)
	data.Set("username", username)
	data.Set("link", link)
	res, err := http.PostForm(YO_API+"/yo/", data)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return errors.New("Received response with non 201 status code.")
	}

	return nil
}
