package agent

import (
    "os"
    "sync"
    "strconv"
    "errors"
    "github.com/abemedia/push-deploy/models"
)

func deploy(p *models.Project, h map[string]string) (errs []error) {
    os.Chdir("compiled")
    var wg sync.WaitGroup
    for i, d := range p.Deploy {
        wg.Add(1)
        go func(i int, d map[string]string, errs []error) {
            defer wg.Done()
            d["log"] = "../logs/deploy-" + strconv.Itoa(i) +".log"
            /*
            switch d.(type) {
                case models.DeployGit: deployGit(deployLog, d.(models.DeployGit), h)
                case models.DeployS3: deployS3(deployLog, d.(models.DeployS3), h)
                default: fmt.Println("Error!")
            }
            */
            var err error
            switch d["type"] {
                case "git": err = deployGit(d, h)
                case "s3": err = deployS3(d, h)
                case "ftp": err = deployFTP(d, h)
                case "local": err = deployLocal(d, h)
                default: err = errors.New("Unknown deployment type.")
            }
            if err != nil {
                errs = append(errs, err)
            }
            /*
            // merge deploy logs into main deploy log file
            os.Create("logs/deploy.log")
            log, _ := ioutil.ReadFile(d["log"])
            logfile.Write(log)
            os.Remove(d["log"])
            */
        }(i, d, errs)
    }
    wg.Wait()
    
    return errs
}
