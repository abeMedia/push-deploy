package agent

import (
    "os"
    "sync"
    "strconv"
    "github.com/abemedia/push-deploy/models"
)

func deploy(p *models.Project, h map[string]string) bool {
    os.Chdir("compiled")
    var wg sync.WaitGroup
    success := true
    for i, d := range p.Deploy {
        wg.Add(1)
        go func(i int, d map[string]string) {
            defer wg.Done()
            d["log"] = "../logs/deploy-" + strconv.Itoa(i) +".log"
            /*
            switch d.(type) {
                case models.DeployGit: deployGit(deployLog, d.(models.DeployGit), h)
                case models.DeployS3: deployS3(deployLog, d.(models.DeployS3), h)
                default: fmt.Println("Error!")
            }
            */
            var ok bool
            switch d["type"] {
                case "git": ok = deployGit(d, h)
                case "s3": ok = deployS3(d, h)
                case "ftp": ok = deployFTP(d, h)
                case "local": ok = deployLocal(d, h)
                default: ok = false
            }
            if !ok {
                success = false
            }
            /*
            // merge deploy logs into main deploy log file
            os.Create("logs/deploy.log")
            log, _ := ioutil.ReadFile(d["log"])
            logfile.Write(log)
            os.Remove(d["log"])
            */
        }(i, d)
    }
    wg.Wait()
    
    return success
}
