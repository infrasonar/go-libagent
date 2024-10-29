package libagent

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
)

type Asset struct {
	Id        int `json:"assetId"`
	Name      string
	Kind      string
	assetFile string
	collector *Collector
}

func NewAsset(collector *Collector) *Asset {
	assetIdEnv := os.Getenv("ASSET_ID")
	if assetIdEnv != "" {
		assetId, err := strconv.Atoi(assetIdEnv)
		if err == nil && assetId > 0 {
			log.Println("Using ASSET_ID environment variable for asset announcement")
			return &Asset{
				Id:        assetId,
				collector: collector,
			}
		}
		log.Println("Invalid ASSET_ID environment variable")
	}

	configPath, err := getConfigPath()
	if err != nil {
		// This is fatal, if no ENV is used, we need a place to write an asset Id
		log.Fatal("Error while reading config path (check the CONFIG_PATH environment variable): ", err)
	}
	assetFile := path.Join(configPath, "asset.json")

	_, err = os.Stat(assetFile)
	if err == nil {
		content, err := os.ReadFile(assetFile)
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}
		// Now let's unmarshall the data into `asset`
		var asset Asset
		err = json.Unmarshal(content, &asset)
		if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
		}
		asset.collector = collector
		asset.assetFile = assetFile
		return &asset
	}

	return &Asset{
		Id:        0,
		assetFile: assetFile,
		collector: collector,
	}
}

func (asset *Asset) Announce() {
	if asset.Id == 0 {
		asset.Create()
		return
	}
	h := GetHelper()
	uri := fmt.Sprintf("/asset/%d?fields=name&collectors=key", asset.Id)
	type Tasset struct {
		Name       string `json:"name"`
		Collectors []struct {
			Key string `json:"key"`
		} `json:"collectors"`
	}
	t := Tasset{}
	err := h.Get(uri, &t)
	if err != nil {
		log.Fatal(err)
	}

	asset.Name = t.Name

	found := false
	for _, v := range t.Collectors {
		if v.Key == asset.collector.Key {
			found = true
			break
		}
	}
	if !found {
		uri = fmt.Sprintf("/asset/%d/collector/%s", asset.Id, asset.collector.Key)
		err = h.Post(uri, nil, nil)
		if err != nil {
			log.Printf("Error while assining collector: %s", err)
		}
	}

	log.Printf("Announced asset `%s` (Id: %d)", asset.Name, asset.Id)
}

func (asset *Asset) Create() {
	if asset.assetFile == "" {
		log.Fatal("missing asset file")
	}

	h := GetHelper()

	type TcontainerId struct {
		ContainerId int `json:"containerId"`
	}
	type Tname struct {
		Name string `json:"name"`
	}
	type TassetId struct {
		AssetId int `json:"assetId"`
	}

	tcid := TcontainerId{}

	err := h.Get("/container/id", &tcid)
	if err != nil {
		log.Fatal(err)
	}

	asset.Name = os.Getenv("ASSET_NAME")
	if asset.Name == "" {
		asset.Name, err = fqdn()
		if err != nil {
			log.Fatal(err)
		}
	}

	data := Tname{Name: asset.Name}
	taid := TassetId{}
	uri := fmt.Sprintf("/container/%d/asset", tcid.ContainerId)

	err = h.Post(uri, &taid, data)
	if err != nil {
		log.Fatal(err)
	}

	asset.Id = taid.AssetId

	uri = fmt.Sprintf("/asset/%d/collector/%s", asset.Id, asset.collector.Key)
	err = h.Post(uri, nil, nil)
	if err != nil {
		log.Printf("Error while assining collector: %s", err)
	}

	if asset.Kind != "" {
		type Tkind struct {
			Kind string `json:"kind"`
		}
		data := Tkind{Kind: asset.Kind}
		uri = fmt.Sprintf("/asset/%d/kind", asset.Id)
		err = h.Patch(uri, nil, data)
		if err != nil {
			log.Printf("Error while setting asset kind: %s", err)
		}
	}

	bytes, _ := json.Marshal(taid)
	err = os.WriteFile(asset.assetFile, bytes, 0644)
	if err != nil {
		log.Printf("Error while writing file (%s): %s", asset.assetFile, err)
	}
	log.Printf("Created asset `%s` (Id: %d)", asset.Name, asset.Id)
}
