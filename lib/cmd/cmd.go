package cmd

import (
    "fmt"
    "os"
    "log"
    "path"
    "path/filepath"
    "os/exec"
    "strings"
)

type command struct {
    name string
    args []string
}

type commands struct {
    commands []command
    logfile string
}

func New(logfile string) commands {
    var c commands
    c.logfile, _ = filepath.Abs(logfile)
    return c
}

// add a shell command to the queue
func (c *commands) Add(name string, args...string) {
    c.commands = append(c.commands, command{name,args})
}

// run queued shell commands
func (c *commands) Run() bool {
    // create project directory if it doesn't exist
    logdir := path.Dir(c.logfile)
    if _, err := os.Stat(logdir); os.IsNotExist(err) {
        os.MkdirAll(logdir, 0755)
    }
    
    // open logfile (create if necessary)
    f, err := os.OpenFile(c.logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        fmt.Println(err)
        return false
    }
    defer f.Close()
    log.SetOutput(f)
    
    // loop through commands & execute
    for _, cmd := range c.commands {
        log.Println("$ ", cmd.name, strings.Join(cmd.args, " "))
        ouput, err := exec.Command(cmd.name, cmd.args...).CombinedOutput()
        log.Println(string(ouput))
        if err != nil {
            return false
        }
    }
    
    // clear command queue
    c.commands = nil
    
    return true
}