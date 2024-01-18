package redis

type Config struct {
	Admin   interface{}   `json:"admin"`
	Storage interface{}   `json:"storage"`
	Logging interface{}   `json:"logging"`
	Adapter AdapterConfig `json:"adapter"`
}

type AdapterConfig struct {
	Prefix           string `json:"prefix"`
	Address          string `json:"address"`
	Password         string `json:"password"`
	Database         int    `json:"database"`
	UpdateInterval   string `json:"updateInterval"`
	SubscribeUpdates string `json:"subscribeUpdates"`
}
