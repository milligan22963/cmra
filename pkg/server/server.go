// Package server is for all server related items
package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

// HTTPResponse is a structure defining what a response should look like
type HTTPResponse struct {
	Code    int    `json:"-"`
	Message string `json:"message,omitempty"`
}

// GenerateHomePage generates the home page for this site
func GenerateHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func setupWebserver(siteConfig *config.SiteConfiguration) {
	httpServerDone := &sync.WaitGroup{}
	router := mux.NewRouter().StrictSlash(true)

	serverPort := viper.GetInt(config.WebServerPort)
	serverAddress := viper.GetString(config.WebServerAddress)
	router.HandleFunc("/", server.GenerateHomePage)
	server := &http.Server{Addr: serverAddress + ":" + strconv.Itoa(serverPort), Handler: router}

	httpServerDone.Add(1) // Add before our go routine
	go func() {
		defer httpServerDone.Done()

		if err := server.ListenAndServe(); err != nil {
			logrus.Errorf("http listen error: %v", err)
		}
	}()

	<-siteConfig.AppActive

	ctx, cancel := context.WithTimeout(context.Background(), httpWait)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("server shutdown error: %v", err)
	}

	// wait for the server func to finish
	httpServerDone.Wait()
}

// ServerInstance is an instance of server
type ServerInstance struct {
	ServerPort int
}

func (server *ServerInstance) Run() {
	// server up the world
	go setupWebserver()

	quitReason := cmd.process(siteConfig)

	siteConfig.AppActive <- struct{}{}

	_ = siteConfig.Database.Close()

}
