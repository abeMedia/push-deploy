package agent

import (
   // "os"
   // "log"
   // "github.com/abemedia/push-deploy/models"
)

const (
    StatusError   int = -1
    StatusEmpty   int = 0
    StatusOK      int = 1
    StatusBuild   int = 2
    StatusDeploy  int = 3
)

var statusText = map[int]string{
    StatusError:  "Error!",
    StatusEmpty:  "No build",
    StatusOK:     "Complete",
    StatusBuild:  "Building...",
    StatusDeploy: "Deploying...",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
    return statusText[code]
}


/*


type statusType struct {
    ID int
    project *models.Project
    log *log.Logger
}

func (s *statusType) update(status int) {
    s.log.Printf("%s: %s", s.project.Name, StatusText(status))
    
    // update the status in the database
    models.Builds.Update(s.ID, map[string]int{"status": status})
    
    // if a channel exists (e.g. websocket is open) write the status to it
    ///if _, ok := Status[p.ID]; ok {
    ///    Status[p.ID] <- status 
    ///}
}

var Status map[int]chan int
func updateStatus(p models.Project, status int) {
    // update the status in the database
    models.Projects.Update(p.ID, map[string]int{"status": status})
    
    logger.Printf("%s: %s", p.Name, StatusText(status))
    
    // if a channel exists (e.g. websocket is open) write the status to it
    if _, ok := Status[p.ID]; ok {
        Status[p.ID] <- status 
    }
}

func GetStatusChannel(project_id int) chan int {
    if _, ok := Status[project_id]; ok {
        return Status[project_id]
    }
    Status[project_id] = make(chan int, 10)
    return Status[project_id]
}
*/