package main

import (
	"fmt"
	"os"
	"time"

	"github.com/libgolang/config"
	"github.com/libgolang/log"
	"gopkg.in/resty.v1"
)

type volumesFlag []string

func (i *volumesFlag) ToString() string {
	return ""
}

func (i *volumesFlag) FromString(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	tokenPtr    = config.String("token", "", "Linode Bearer Token")
	namePtr     = config.String("name", "", "Container Name")
	hostPtr     = config.String("host", getHostName(), "Hostname to attach volume to")
	hookTypePtr = config.String("hook", "", "Hook Type: pre | post")
	volumes     volumesFlag
)

func main() {
	resty.SetDebug(true)
	_ = os.Setenv("LOG_CONFIG", "config.properties")
	log.LoadLogProperties()

	//config.String("config", "one-linode.conf", "Path to config file")
	config.Var(&volumes, "volume", "Volume to attach. Takes multiple volumes. E.g: --volume vol1 --volume vol2")
	config.Parse()

	if *tokenPtr == "" {
		fmt.Printf("###################################\n")
		fmt.Printf("--token or $TOKEN confg is required\n")
		fmt.Printf("###################################\n")
	} else if *hookTypePtr == "pre" {
		preHook()
	} else if *hookTypePtr == "post" {
		postHook()
	} else {
		fmt.Printf("##################################################\n")
		fmt.Printf("--hook flag is required. Possible values: pre|post\n")
		fmt.Printf("##################################################\n")
	}
}

func preHook() {
	for _, volumeName := range volumes {
		if err := attachLinode(*hostPtr, volumeName); err != nil {
			os.Exit(1) // error exit
		}
	}
}

func postHook() {
}

func attachLinode(linodeName string, volumeName string) error {
	linodeID, err := getLinodeIDByName(linodeName)
	if err != nil {
		err = fmt.Errorf("Unable to get Linode ID by name(%s): %s", linodeName, err)
		log.Error("%s", err)
		return err
	}

	volumeID, err := getVolumeIDByName(volumeName)
	if err != nil {
		err = fmt.Errorf("Unable to get Volume ID by name(%s)", volumeName)
		log.Error("%s", err)
		return err
	}

	// detach
	log.Info("Calling detach on volume %d", volumeID)
	detachURL := fmt.Sprintf("https://api.linode.com/v4/volumes/%d/detach", volumeID)
	if _, err := Post(detachURL, nil, nil); err != nil {
		log.Warn("Detaching request returned error")
	}
	// wait for deatch request to finish
	i := 0
	for {
		duration := time.Second * 5
		log.Info("Wait for deatch request %s", duration)
		time.Sleep(duration) // sleep 5 seconds

		it, err := Get(fmt.Sprintf("https://api.linode.com/v4/volumes/%d", volumeID), &Volume{})
		if err != nil {
			log.Error("Detach Wait request failed")
		}
		vol := it.(*Volume)
		if vol.LinodeID == 0 {
			log.Info("Node detached stop the wait")
			break
		}
		if i >= 20 {
			break
		}
		i++
	}

	// attach
	log.Info("Calling attach on volume %d and node %d", volumeID, linodeID)
	url := fmt.Sprintf("https://api.linode.com/v4/volumes/%d/attach", volumeID)
	body := AttachRequest{LinodeID: &linodeID}
	if _, err := Post(url, body, nil); err != nil {
		err = fmt.Errorf("unable to attach volume: %s", err)
		log.Error("%s", err)
		return err
	}
	return nil
}

// getLinodeIDByName resturns the id of the linode given the name or returns empty
// string if not found
func getLinodeIDByName(linodeName string) (int, error) {
	pages := 1
	for page := 1; page <= pages; page++ {
		url := fmt.Sprintf("https://api.linode.com/v4/linode/instances?page=%d", page)
		//var err error
		it, err := Get(url, &ListNodeResponse{})
		if err != nil {
			return 0, err
		}
		resp, ok := it.(*ListNodeResponse)
		if !ok {
			return 0, fmt.Errorf("Error casting to ListNodeReponse")
		}
		pages = resp.Pages
		for _, n := range resp.Data {
			if n.Label == linodeName {
				return n.ID, nil
			}
		}
	}
	return 0, fmt.Errorf("Not Found")
}

func getVolumeIDByName(volumeName string) (int, error) {
	pages := 1
	for page := 1; page <= pages; page++ {
		url := fmt.Sprintf("https://api.linode.com/v4/volumes?page=%d", page)
		it, err := Get(url, &ListVolumeResponse{})
		if err != nil {
			return 0, err
		}
		resp, ok := it.(*ListVolumeResponse)
		if !ok {
			return 0, fmt.Errorf("Error casting to ListVolumeReponse")
		}
		pages = resp.Pages
		for _, n := range resp.Data {
			if n.Label == volumeName {
				return n.ID, nil
			}
		}
	}
	return 0, fmt.Errorf("Not Found")
}

// Get REST GET request
func Get(url string, res interface{}) (interface{}, error) {
	log.Debug("GET %s token: %s", url, *tokenPtr)
	r := resty.R()
	if res != nil {
		r.SetResult(res)
	}
	r.SetHeader("Authorization", fmt.Sprintf("Bearer %s", *tokenPtr))
	resp, err := r.Get(url)
	if err == nil && resp.StatusCode() != 200 {
		return nil, fmt.Errorf("GET Request returned error %d: ", resp.StatusCode())
	}
	return resp.Result(), err
}

// Post REST POST request
func Post(url string, req interface{}, res interface{}) (interface{}, error) {
	log.Debug("POST %s", url)

	r := resty.R()
	if req != nil {
		r.SetBody(req)
	}

	if res != nil {
		r.SetResult(res)
	}

	r.SetHeader("Authorization", fmt.Sprintf("Bearer %s", *tokenPtr))
	resp, err := r.Post(url)
	if err == nil && resp.StatusCode() != 200 {
		return nil, fmt.Errorf("GET Request returned error %d: ", resp.StatusCode())
	}
	return resp.Result(), err
}

// ListNodeResponse list node response
type ListNodeResponse struct {
	Data    []Node `json:"data"`
	Page    int    `json:"page"`    // "page": 1,
	Pages   int    `json:"pages"`   // "pages": 1,
	Results int    `json:"results"` // "results": 1
}

// ListVolumeResponse list volume response
type ListVolumeResponse struct {
	Data    []Volume `json:"data"`
	Page    int      `json:"page"`    // "page": 1,
	Pages   int      `json:"pages"`   // "pages": 1,
	Results int      `json:"results"` // "results": 1
}

// Node node
type Node struct {
	ID     int    `json:"id"`     //"id": 123,
	Label  string `json:"label"`  //"label": "linode123",
	Region string `json:"region"` //"region": "us-east",
	//"image": "linode/debian9",
	//"type": "g6-standard-2",
	//"group": "Linode-Group",
	//"status": "running",
	//"hypervisor": "kvm",
	//"created": "2018-01-01T00:01:01",
	//"updated": "2018-01-01T00:01:01",
	//...
	//...
	//...
}

// Volume volume
type Volume struct {
	ID             int    `json:"id"`              // "id": 12345,
	Label          string `json:"label"`           // "label": "my-volume",
	FilesystemPath string `json:"filesystem_path"` // "filesystem_path": "/dev/disk/by-id/scsi-0Linode_Volume_my-volume",
	LinodeID       int    `json:"linode_id"`       // "linode_id": 12346,
	Region         string `json:"region"`          // "region": "us-east",
	// "status": "active",
	// "size": 30,
	// "created": "2018-01-01T00:01:01",
	// "updated": "2018-01-01T00:01:01"
}

// AttachRequest linode Attach Request
type AttachRequest struct {
	LinodeID *int    `json:"linode_id"`
	ConfigID *string `json:"config_id"`
}

func getHostName() string {
	h, _ := os.Hostname()
	return h
}
