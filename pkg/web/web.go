// package web implements the web interface for this app
package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/milligan22963/afmlog"
	"github.com/milligan22963/cmra/config"
)

const (
	httpWait = 5 * time.Second
)

// WebServer is the web serving object for this app
type WebServer struct {
	logger *afmlog.Log
}

func (webserver *WebServer) SetupWebserver(appConfig *config.AppConfiguration) {
	webserver.logger = appConfig.GetLogger()
	httpServerDone := &sync.WaitGroup{}
	router := mux.NewRouter().StrictSlash(true)

	serverPort := appConfig.CameraConfiguration.WebServerSettings.Port
	serverAddress := appConfig.CameraConfiguration.WebServerSettings.Host
	router.HandleFunc("/", webserver.GenerateHomePage)
	server := &http.Server{Addr: serverAddress + ":" + strconv.Itoa(serverPort), Handler: router}

	webserver.logger.Informationf("Address: %v, Port: %d", serverAddress, serverPort)
	httpServerDone.Add(1) // Add before our go routine
	go func() {
		defer httpServerDone.Done()

		if err := server.ListenAndServe(); err != nil {
			webserver.logger.Errorf("http listen error: %v", err)
		}
	}()

	<-appConfig.AppActive

	ctx, cancel := context.WithTimeout(context.Background(), httpWait)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		webserver.logger.Errorf("server shutdown error: %v", err)
	}

	// wait for the server func to finish
	httpServerDone.Wait()
}

// GenerateHomePage generates the home page for this site
func (webserver *WebServer) GenerateHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}
