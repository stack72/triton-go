package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/joyent/triton-go/client"
)

type VolumesClient struct {
	client *client.Client
}

type Volume struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	OwnerID        string    `json:"owner_uuid"`
	Type           string    `json:"type"`
	State          string    `json:"state"`
	Networks       []string  `json:"networks"`
	Refs           []string  `json:"refs"`
	Size           int64     `json:"size"`
	FileSystemPath string    `json:"filesystem_path"`
	Created        time.Time `json:"create_timestamp"`
}

type ListVolumeInput struct {
	Name  string
	Size  int64
	State string
	Type  string
}

func (c *VolumesClient) List(ctx context.Context, input *ListVolumeInput) ([]*Volume, error) {
	path := fmt.Sprintf("/%s/volumes", c.client.AccountName)

	query := &url.Values{}
	if input.Name != "" {
		query.Set("name", input.Name)
	}

	if input.Size != 0 {
		query.Set("size", fmt.Sprintf("%d", input.Size))
	}

	if input.State != "" {
		query.Set("state", input.State)
	}

	if input.Type != "" {
		query.Set("type", input.Type)
	}

	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
		Query:  query,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing List request: {{err}}", err)
	}

	var result []*Volume
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding List response: {{err}}", err)
	}

	return result, nil
}

type GetVolumeInput struct {
	ID string
}

func (c *VolumesClient) Get(ctx context.Context, input *GetVolumeInput) (*Volume, error) {
	path := fmt.Sprintf("/%s/volumes/%s", c.client.AccountName, input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing Get request: {{err}}", err)
	}

	var result *Volume
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding Get response: {{err}}", err)
	}

	return result, nil
}

type CreateVolumeInput struct {
	Name     string   `json:"name,omitempty"`
	Size     int64    `json:"size,omitempty"`
	Type     string   `json:"type"`
	Networks []string `json:"networks"`
}

func (c *VolumesClient) Create(ctx context.Context, input *CreateVolumeInput) (*Volume, error) {
	path := fmt.Sprintf("/%s/volumes", c.client.AccountName)

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Body:   input,
	}

	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing Create request: {{err}}", err)
	}

	var result *Volume
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding Create response: {{err}}", err)
	}

	return result, nil
}

type DeleteVolumeInput struct {
	ID string
}

func (c *VolumesClient) Delete(ctx context.Context, input *DeleteVolumeInput) error {
	path := fmt.Sprintf("/%s/volumes/%s", c.client.AccountName, input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing Delete request: {{err}}", err)
	}

	return nil
}

type UpdateVolumeInput struct {
	ID   string
	Name string
}

func (c *VolumesClient) Update(ctx context.Context, input *UpdateVolumeInput) error {
	path := fmt.Sprintf("/%s/volumes/%s", c.client.AccountName, input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Body:   input.Name,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing Update request: {{err}}", err)
	}

	return nil
}
