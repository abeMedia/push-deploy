package agent

import (
    "os"
    "github.com/abemedia/push-deploy/lib/cmd"
    "github.com/abemedia/push-deploy/models"
)

func build(p *models.Project) error {
    command := cmd.New("logs/build.log")
    
    // change working directory to repo
    os.Mkdir("source", os.ModePerm)
    os.Chdir("source")
    
    // init git repo
    if _, err := os.Stat(".git"); os.IsNotExist(err) {
        command.Add("git", "init")
        command.Add("git", "remote", "add", "origin", p.Repo)
    }
    
    // pull project code
    command.Add("git", "fetch", "--no-tags", "--depth=1", "origin", "+refs/heads/" + p.Branch)
    command.Add("git", "reset", "--hard", "origin/" + p.Branch)
    if err := command.Run(); err != nil {
        return err
    }
    
    buildJekyll()
    
    // change working directory back to project root
    os.Chdir("../")
    
    return nil
}