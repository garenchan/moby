package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
)

// GetDownloadBandwidth returns max download bandwidth about the docker server.
func (cli *Client) GetDownloadBandwidth(ctx context.Context) (int64, error) {
	var bandwidth int64

	serverResp, err := cli.get(ctx, "/bandwidth/download", nil, nil)
	defer ensureReaderClosed(serverResp)
	if err != nil {
		return bandwidth, nil
	}

	if err := json.NewDecoder(serverResp.body).Decode(&bandwidth); err != nil {
		return bandwidth, fmt.Errorf("Error get max download bandwidth: %v", err)
	}

	return bandwidth, nil
}

// SetDownloadBandwidth sets max download bandwidth about the docker server.
func (cli *Client) SetDownloadBandwidth(ctx context.Context, bandwidth int64, persistent bool) error {
	bd := types.Bandwidth{Value: bandwidth, Persistent: persistent}

	resp, err := cli.post(ctx, "/bandwidth/download", nil, bd, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return err
	}

	return nil
}

// GetUploadBandwidth returns max upload bandwidth about the docker server.
func (cli *Client) GetUploadBandwidth(ctx context.Context) (int64, error) {
	var bandwidth int64

	serverResp, err := cli.get(ctx, "/bandwidth/upload", nil, nil)
	defer ensureReaderClosed(serverResp)
	if err != nil {
		return bandwidth, nil
	}

	if err := json.NewDecoder(serverResp.body).Decode(&bandwidth); err != nil {
		return bandwidth, fmt.Errorf("Error get max upload bandwidth: %v", err)
	}

	return bandwidth, nil
}

// SetUploadBandwidth sets max upload bandwidth about the docker server.
func (cli *Client) SetUploadBandwidth(ctx context.Context, bandwidth int64, persistent bool) error {
	bd := types.Bandwidth{Value: bandwidth, Persistent: persistent}

	resp, err := cli.post(ctx, "/bandwidth/upload", nil, bd, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return err
	}

	return nil
}
