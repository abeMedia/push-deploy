package agent

import (
    "os"
    "io/ioutil"
    "fmt"
    "log"
    "time"
    "github.com/abemedia/push-deploy/models"
    . "github.com/abemedia/push-deploy/lib/config"
)

func Run(h map[string]string, p models.Project) {
    logger := log.New(os.Stdout, "["+p.Name+"] ", log.Ldate|log.Ltime)
    last_build, err := models.Builds.Get(map[string]interface{}{"project_id": p.ID})
    
    if queue, _ := models.Queues.Get(map[string]interface{}{"project_id": p.ID}); queue.ID > 0 {
        // queued & running - exit silently
        if last_build.Status >= StatusBuild {
            return
        }
        
        // queued but not running. remove from queue
        models.Queues.Delete(queue.ID)
    }
    
    // build already running - add to queue and exit
    if last_build.Status >= StatusBuild {
        models.Queues.Add(map[string]interface{}{"project_id": p.ID})
        logger.Print("Queued")
        return
    }
    
    // create new build in DB
    build_info := &models.Build{
        ProjectID: p.ID,
        Start: time.Now(),
        Author: h["author"],
        Message: h["message"],
    }
    id, err := models.Builds.Add(build_info)
    if err != nil {
        logger.Println(err)
        return
    }
    
    // set status to building
    models.Builds.Update(id, map[string]interface{}{"status": StatusBuild})
    logger.Print(StatusText(StatusBuild))
    
    // get & switch to to project directory - create if it doesn't exist
    dir := fmt.Sprintf("%s/%d", Config.CachePath, p.ID)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        os.MkdirAll(dir, 0777)
    }
    os.Chdir(dir)
    logdir := dir + "/logs"
    
    if !build(&p) {
        models.Builds.Update(id, map[string]interface{}{"status": StatusError})
        logger.Print(StatusText(StatusError))
        return
    }
    
    // set status to deploying
    models.Builds.Update(id, map[string]interface{}{"status": StatusDeploy})
    logger.Print(StatusText(StatusDeploy))
    
    if !deploy(&p, h) {
        models.Builds.Update(id, map[string]interface{}{"status": StatusError})
        logger.Print(StatusText(StatusError))
        return
    }
    
    // set status to live & add build info to database
    build_info.Finish = time.Now()
    build_info.Status = StatusOK
    build_info.Logs = make(map[string]string)
    files, _ := ioutil.ReadDir(logdir)
    for _, f := range files {
        log, _ := ioutil.ReadFile(logdir + "/" + f.Name())
        build_info.Logs[f.Name()] = string(log)
    }
    models.Builds.Update(id, build_info)
    
    // cleanup
    os.RemoveAll(dir + "/compiled")
    os.RemoveAll(logdir)
    
    logger.Printf(StatusText(StatusOK))
    
    // if this project is in the queue, build it again
    if queue, _ := models.Queues.Get(map[string]interface{}{"project_id": p.ID}); queue.ID > 0 {
        project, _ := models.Projects.Get(map[string]interface{}{"id": p.ID})
        models.Queues.Delete(queue.ID)
        Run(h, project)
    }
}