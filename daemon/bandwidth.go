package daemon

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/docker/docker/pkg/ioutils"
	metrics "github.com/docker/go-metrics"
	"github.com/sirupsen/logrus"
)

const (
	// downBwConfig download bandwidth config
	downBwConfig = "max-download-bandwidth"
	// upBwConfig upload bandwidth config
	upBwConfig = "max-upload-bandwidth"
)

// GetDownloadBandwidth returns max download bandwidth about the daemon.
func (d *Daemon) GetDownloadBandwidth() int64 {
	defer metrics.StartTimer(hostInfoFunctions.WithValues("get_download_bandwidth"))()
	return d.imageService.GetDownloadBandwidth()
}

// SetDownloadBandwidth sets max download bandwidth about the daemon.
func (d *Daemon) SetDownloadBandwidth(bandwidth int64, persistent bool) error {
	logrus.Infof("Set max download bandwidth to %d, persistent=%t", bandwidth, persistent)
	defer metrics.StartTimer(hostInfoFunctions.WithValues("set_download_bandwidth"))()

	if persistent {
		if err := d.persistBandwidthConfig(downBwConfig, bandwidth); err != nil {
			return err
		}
	}

	d.imageService.UpdateDownloadBandwidth(&bandwidth)
	return nil
}

// GetUploadBandwidth returns max upload bandwidth about the daemon.
func (d *Daemon) GetUploadBandwidth() int64 {
	defer metrics.StartTimer(hostInfoFunctions.WithValues("get_upload_bandwidth"))()
	return d.imageService.GetUploadBandwidth()
}

// SetUploadBandwidth sets max upload bandwidth about the daemon.
func (d *Daemon) SetUploadBandwidth(bandwidth int64, persistent bool) error {
	logrus.Infof("Set max upload bandwidth to %d, persistent=%t", bandwidth, persistent)
	defer metrics.StartTimer(hostInfoFunctions.WithValues("set_upload_bandwidth"))()

	if persistent {
		if err := d.persistBandwidthConfig(upBwConfig, bandwidth); err != nil {
			return err
		}
	}

	d.imageService.UpdateUploadBandwidth(&bandwidth)
	return nil
}

// persistBandwidthConfig persist bandwidth related config to file.
func (d *Daemon) persistBandwidthConfig(name string, value interface{}) error {
	d.configStore.Lock()
	defer d.configStore.Unlock()

	configFile := d.configStore.ConfigFile
	var jsonConfig map[string]interface{}

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		// config file maybe not exist
		if !os.IsNotExist(err) {
			return err
		}
		jsonConfig = make(map[string]interface{})
	} else {
		reader := bytes.NewReader(b)
		if err := json.NewDecoder(reader).Decode(&jsonConfig); err != nil {
			return nil
		}
	}

	// change config
	jsonConfig[name] = value

	if b, err = json.MarshalIndent(jsonConfig, "", "    "); err != nil {
		return err
	}

	if err := ioutils.AtomicWriteFile(configFile, b, os.FileMode(0600)); err != nil {
		return err
	}

	return nil
}
