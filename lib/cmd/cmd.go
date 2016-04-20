package cmd

import (
    "os"
    "log"
    "path"
    "path/filepath"
    "os/exec"
    "strings"
    "time"
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
func (c *commands) Run() error {
    // create project directory if it doesn't exist
    logdir := path.Dir(c.logfile)
    if _, err := os.Stat(logdir); os.IsNotExist(err) {
        os.MkdirAll(logdir, 0755)
    }
    
    // open logfile (create if necessary)
    f, err := os.OpenFile(c.logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        return err
    }
    defer f.Close()
    log.SetOutput(f)
    
    // loop through commands & execute
    for _, cmd := range c.commands {
        log.Println("$", cmd.name, strings.Join(cmd.args, " "))
        start := time.Now()
        ouput, err := exec.Command(cmd.name, cmd.args...).CombinedOutput()
        end := time.Now()
        log.Print(string(ouput))
        log.Println("Completed in", end.Sub(start).String(), "\n")
        if err != nil {
            return err
        }
    }
    
    // clear command queue
    c.commands = nil
    
    return nil
}