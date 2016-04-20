package config

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "flag"
    "path/filepath"
    "os"
    "path"
)

type Configuration struct {
    Host        string
    Port        int
    CachePath   string `yaml:"cache_path"`
    CacheTime   string `yaml:"cache_time"`
    DB struct {
		    Type        string
		    Host        string
		    Database    string
		    User        string
		    Password    string
    }
    DashboardPath   string `yaml:"dashboard_path"`
}

var Config Configuration

func init() {
    // get config file path from command line flag
    configFile := flag.String("c", "./push-deploy.conf", "path to config file")
    flag.Parse()
    
    // open config file
    file, err := ioutil.ReadFile(*configFile)
    if err != nil {
        panic(err.Error())
    }
    
    // parse yaml into struct
    err = yaml.Unmarshal(file, &Config)
    if err != nil {
        panic(err.Error())
    }
    
    if Config.CachePath == "" {
        Config.CachePath = path.Join(os.TempDir(), "./jekyll-deploy")
    } else {
        Config.CachePath, err = filepath.Abs(Config.CachePath)
        if err != nil {
            panic(err.Error())
        }
    }
    
    if Config.DashboardPath == "" {
        Config.DashboardPath, _ = filepath.Abs("./static")
    } else {
        Config.DashboardPath, err = filepath.Abs(Config.DashboardPath)
        if err != nil {
            panic(err.Error())
        }
    }
}