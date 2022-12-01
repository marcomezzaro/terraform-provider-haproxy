package haproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/matthisholleville/terraform-provider-haproxy/internal/haproxy/models"
)

func (c *Client) GetBind(bind models.Bind, parent_name string) (*models.Bind, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/binds/%s?frontend=%s&parent_type=frontend&parent_name=%s", c.base_url, bind.Name, parent_name, parent_name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetBind{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (c *Client) CreateBind(transactionId string, bind models.Bind, parent_name string) (*models.Bind, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/binds?transaction_id=%s&frontend=%s&parent_type=frontend&parent_name=%s", c.base_url, transactionId, parent_name, parent_name)
	bodyStr, _ := json.Marshal(bind)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Bind{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateBind(transactionId string, bind models.Bind, parent_name string) (*models.Bind, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/binds/%s?transaction_id=%s&frontend=%s&parent_type=frontend&parent_name=%s", c.base_url, bind.Name, transactionId, parent_name, parent_name)
	bodyStr, _ := json.Marshal(bind)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Bind{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteBind(transactionId string, bind models.Bind, parent_name string) error {
	url := fmt.Sprintf("%s/services/haproxy/configuration/binds/%s?transaction_id=%s&frontend=%s&parent_type=frontend&parent_name=%s", c.base_url, bind.Name, transactionId, parent_name, parent_name)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}
