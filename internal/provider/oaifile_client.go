package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type OAIFileClient struct {
	Host string
}

type OAIFile struct {
	ID   int   `json:"id"`
	FileID string `json:"file_id"`
	FilePath string `json:"filepath"`
	Name string `json:"name"`
}

type OAIFileRequest struct {
	FilePath 	string `json:"filepath"`
	Name 			string `json:"name"`
}

func (c *OAIFileClient) do(req *http.Request, v interface{}) error {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if v != nil {
		dec := json.NewDecoder(res.Body)
		return dec.Decode(v)
	}
	return nil
}

func (c *OAIFileClient) GetOAIFiles(ctx context.Context) ([]OAIFile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Host+"/data", nil)
	if err != nil {
		return nil, err
	}

	var oaifiles []OAIFile
	if err := c.do(req, oaifiles); err != nil {
		return nil, err
	}

	return oaifiles, nil
}

func (c *OAIFileClient) GetOAIFile(ctx context.Context, id int) (*OAIFile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Host+"/data/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	var oaifile OAIFile
	if err := c.do(req, &oaifile); err != nil {
		return nil, err
	}

	return &oaifile, nil
}

func (c *OAIFileClient) CreateOAIFile(ctx context.Context, oaifile OAIFileRequest) (*OAIFile, error) {
	jsonData, err := json.Marshal(oaifile)
	if err != nil {
		return nil, err
	}
	tflog.Info(ctx, "json data:", map[string]any{"request": string(jsonData)})
	bodyReader := bytes.NewBuffer(jsonData)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Host+"/data", bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var createdOAIFile OAIFile
	if err := c.do(req, &createdOAIFile); err != nil {
		return nil, err
	}
	
	jsonData, err = json.Marshal(createdOAIFile)
	if err != nil {
		return nil, err
	}
	tflog.Info(ctx, "json output:", map[string]any{"data": string(jsonData)})

	return &createdOAIFile, nil
}

func (c *OAIFileClient) UpdateOAIFile(ctx context.Context, oaifile OAIFile) (*OAIFile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.Host+"/data/"+strconv.Itoa(oaifile.ID), nil)
	if err != nil {
		return nil, err
	}

	var updatedOAIFile OAIFile
	if err := c.do(req, &updatedOAIFile); err != nil {
		return nil, err
	}

	return &updatedOAIFile, nil
}

func (c *OAIFileClient) DeleteOAIFile(ctx context.Context, id int) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.Host+"/data/"+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}

	if err := c.do(req, nil); err != nil {
		return err
	}

	return nil
}
