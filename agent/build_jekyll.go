package agent

import (
    "os"
    "github.com/abemedia/push-deploy/lib/cmd"
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

func buildJekyll() bool {
    
    // update config file
    jekyllConfig()
    
    command := cmd.New("../logs/build.log")
    
    // build jekyll site
    if _, err := os.Stat("Gemfile"); err == nil {
        command.Add("bundle", "install", "--deployment", "--clean", "--path", "../vendor")
        command.Add("bundle", "exec", "jekyll", "build", "--destination", "../compiled")
    } else {
        command.Add("jekyll", "build", "--destination", "../compiled")
    }
    if !command.Run() {
        return false
    }
    
    return true
}

func jekyllConfig() error {
    // open config file
    file, err := ioutil.ReadFile("_config.yml")
    if err != nil {
        return err
    }
    
    // parse jekyll config
    var c map[string]interface{}
    err = yaml.Unmarshal(file, &c)
    if err != nil {
        return err
    }
    
    // add default exclude params
    if _, ok := c["exclude"]; !ok {
        c["exclude"] = []string{}
    }
    c["exclude"] = append(c["exclude"].([]string), []string{"Gemfile","Gemfile.lock","bower.json"}...)
    
    // encode back to yaml
    d, err := yaml.Marshal(&c)
    if err != nil {
        return err
    }
    
    // write new jekyll config file
    err = ioutil.WriteFile("_config.yml", d, 0644)
    if err != nil {
        return err
    }
    
    return nil
}