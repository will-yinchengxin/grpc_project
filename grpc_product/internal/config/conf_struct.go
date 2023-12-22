package config

type AppConfig struct {
	AppInfo          AppInfo          `json:"app"`
	DBConfig         DBConfig         `json:"db"`
	ConsulConfig     ConsulConfig     `json:"consul"`
	JWTConfig        JWTConfig        `json:"jwt"`
	ProductWebConfig ProductWebConfig `json:"productWeb"`
}

type AppInfo struct {
	SrvName string   `json:"srvName"`
	SrvTag  []string `json:"srvTag"`
}
type NacosConfig struct {
	Host      string `json:"host"`
	Port      uint64 `json:"port"`
	NameSpace string `json:"namespace"`
	DataId    string `json:"dataid"`
	Group     string `json:"group"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"dbName"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type JWTConfig struct {
	SingingKey string `json:"key" json:"key"`
}

type ProductWebConfig struct {
	SrvName string   `json:"srvName"`
	Host    string   `json:"host"`
	Port    int      `json:"port"`
	Tags    []string `json:"tags"`
}
