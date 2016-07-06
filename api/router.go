package api

import (
    "net/http"
    "github.com/gorilla/mux"
)

/*
// for writing one(two(three))) like use(three, two, one)
func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}
*/


func init() {
    var router = mux.NewRouter().PathPrefix("/api").Subrouter()
    
    // projects
    router.HandleFunc("/projects", auth(Projects)).Methods("GET")
    router.HandleFunc("/projects", auth(ProjectAdd)).Methods("POST")
    router.HandleFunc("/projects/{project_id:[0-9]+}", auth(ProjectView)).Methods("GET")
    router.HandleFunc("/projects/{project_id:[0-9]+}", auth(ProjectUpdate)).Methods("PUT")
    router.HandleFunc("/projects/{project_id:[0-9]+}", auth(ProjectDelete)).Methods("DELETE")
    
    // project builds
    router.HandleFunc("/projects/{project_id:[0-9]+}/builds", auth(Builds)).Methods("GET")
    router.HandleFunc("/projects/{project_id:[0-9]+}/builds", auth(BuildNew)).Methods("POST")
    router.HandleFunc("/projects/{project_id:[0-9]+}/builds/{build_id:[0-9]+}", auth(BuildView)).Methods("GET")
    router.HandleFunc("/projects/{project_id:[0-9]+}/builds/{build_id:[0-9]+}", auth(BuildDelete)).Methods("DELETE")
    
    // status websockets
    router.HandleFunc("/statuses/{user_id:[0-9]+}", ProjectsStatus)
    router.HandleFunc("/status/{id:[0-9]+}", ProjectStatus)
    
    router.HandleFunc("/log/{id:[0-9]+}", auth(Log)).Methods("GET")
    
    // users
    router.HandleFunc("/users", auth(Users)).Methods("GET")
    router.HandleFunc("/users/{id:[0-9]+}", auth(UserView)).Methods("GET")
    router.HandleFunc("/users/{id:[0-9]+}", auth(UserUpdate)).Methods("PUT")
    router.HandleFunc("/users", auth(UserAdd)).Methods("POST")
    router.HandleFunc("/users/{id:[0-9]+}", auth(UserDelete)).Methods("DELETE")
    
    
    router.HandleFunc("/login", Login)
    
    
    //router.HandleFunc("/login", loginHandler)
    //router.HandleFunc("/logout", logoutHandler)
    //router.HandleFunc("/api/logout", api.Logout)
    
    http.Handle("/api/", router)
}
