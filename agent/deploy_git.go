package agent

import (
    "os"
    "fmt"
    "github.com/abemedia/push-deploy/lib/cmd"
)

func deployGit(deploy map[string]string, hook map[string]string) error {
    command := cmd.New(deploy["log"])
    
    // initiate git repo if it isn't already done
    if _, err := os.Stat(".git"); os.IsNotExist(err) {
        command.Add("git", "init")
        command.Add("git", "remote", "add", "origin", deploy["repo"])
        
        // todo: offer option to keep intact history, e.g. performing a git fetch & merge
    }
    
    // add all files, commit & push 
    command.Add("git", "add", "-A")
    command.Add("git", "commit", "-am", hook["message"], fmt.Sprintf(`--author="%s <%s>"`, hook["author"], hook["email"]))
    command.Add("git", "push", "origin", "HEAD:" + deploy["branch"], "-f")
    return command.Run()
}