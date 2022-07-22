package chain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/log"
	"io"
	"net/http"
	"time"
)

type LinkClient struct {
	Config *config.Config
}

// RegisterResource Register chain
func (c *LinkClient) RegisterResource(r ResourceInfo) error {

	old, err := c.LoadRegistryInfoFromChain()

	if err != nil || old.PeerId != c.Config.Identity.PeerID {

		data, err := json.Marshal(r)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", c.Config.LinkApi+"/api/resources", bytes.NewBuffer(data))

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		resp, err := client.Do(req)

		if err != nil {
			return err
		}

		defer resp.Body.Close()

		log.GetLogger().Info("successfully registered virtual machine")
	}
	return nil
}

func (c *LinkClient) RemoveResource(index uint64) error {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/resources/%d", c.Config.LinkApi, index), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil

}

func (c *LinkClient) ModifyResourcePrice(index uint64, unitPrice int64) error {
	return nil
}

func (c *LinkClient) ChangeResourceStatus(index uint64) error {
	return nil
}

func (c *LinkClient) AddResourceDuration(index uint64, duration int) error {
	return nil
}

func (c *LinkClient) Heartbeat(agreementindex uint64) error {
	return nil
}

// LoadRegistryInfoFromChain load registration information from the chain
func (c *LinkClient) LoadRegistryInfoFromChain() (*ResourceInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.Config.LinkApi+"/api/resource/peer?peerId="+c.Config.Identity.PeerID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	//parse body
	var info ResourceInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.GetLogger().Error(err)
		return nil, err
	}
	return &info, nil
}

// LoadKeyFromChain Get the public Yue of the current node from the chain
func (c *LinkClient) LoadKeyFromChain() ([]string, error) {

	info, err := c.LoadRegistryInfoFromChain()
	if err != nil {
		return []string{}, err
	}
	res := []string{info.User}
	return res, nil
}

// Heatbeat ReportStatus Timing the virtual machine to report the service status
func (c *LinkClient) Heatbeat(agreementindex uint64) error {
	return nil
}

func (c *LinkClient) OrderExec(orderIndex uint64) error {
	return nil
}

func (c *LinkClient) CalculateInstanceOverdue(agreementIndex uint64) time.Duration {
	return time.Second
}
