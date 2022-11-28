package haproxy

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/matthisholleville/terraform-provider-haproxy/internal/haproxy/models"
)

func (c *Client) GetBackend(backend models.Backend) (*models.Backend, error) {
	url := c.base_url + "/services/haproxy/configuration/backends/" + backend.Name
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetBackend{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (c *Client) CreateBackend(transactionId string, backend models.Backend) (*models.Backend, error) {
	url := c.base_url + "/services/haproxy/configuration/backends?transaction_id=" + transactionId
	bodyStr, _ := json.Marshal(backend)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Backend{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateBackend(transactionId string, backend models.Backend) (*models.Backend, error) {
	url := c.base_url + "/services/haproxy/configuration/backends/" + backend.Name + "?transaction_id=" + transactionId
	bodyStr, _ := json.Marshal(backend)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Backend{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteBackend(transactionId string, backend models.Backend) error {
	url := c.base_url + "/services/haproxy/configuration/backends/" + backend.Name + "?transaction_id=" + transactionId
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}
