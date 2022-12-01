package haproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/matthisholleville/terraform-provider-haproxy/internal/haproxy/models"
)

func (c *Client) GetServer(server models.Server, parent_name string) (*models.Server, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/servers/%s?backend=%s&parent_type=backend&parent_name=%s", c.base_url, server.Name, parent_name, parent_name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetServer{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (c *Client) CreateServer(transactionId string, server models.Server, parent_name string) (*models.Server, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/servers?transaction_id=%s&backend=%s&parent_type=backend&parent_name=%s", c.base_url, transactionId, parent_name, parent_name)
	bodyStr, _ := json.Marshal(server)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Server{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateServer(transactionId string, server models.Server, parent_name string) (*models.Server, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/servers/%s?transaction_id=%s&backend=%s&parent_type=backend&parent_name=%s", c.base_url, server.Name, transactionId, parent_name, parent_name)
	bodyStr, _ := json.Marshal(server)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Server{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteServer(transactionId string, server models.Server, parent_name string) error {
	url := fmt.Sprintf("%s/services/haproxy/configuration/servers/%s?transaction_id=%s&backend=%s&parent_type=backend&parent_name=%s", c.base_url, server.Name, transactionId, parent_name, parent_name)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}
