package app

import (
	"encoding/json"
	"log"
	"net/http"
	"splash/processing"
	"splash/services"
)

type AnalyticsServer struct {
	config map[string]string
}

func NewAnalyticsServer(config map[string]string) *AnalyticsServer {
	return &AnalyticsServer{
		config: config,
	}
}

func (self *AnalyticsServer) Start() {

	listenAddr := self.config["host"] + ":" + self.config["port"]

	http.HandleFunc(self.config["path"], analyticsHandler(self))

	// Should be last line as it is a blocking.
	log.Fatal("Analytics Server: ", http.ListenAndServe(listenAddr, nil))
}

func analyticsHandler(self *AnalyticsServer) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		self.respond(w)
	}
}

func (*AnalyticsServer) respond(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Locking on Aggregated data
	dailyActiveUsers := <-processing.AggregationQueue

	serviceLocator := services.Locator{}
	logger := serviceLocator.Logger()

	logger.Info("finally", dailyActiveUsers)

	json.NewEncoder(w).Encode(map[string]interface{}{"dailyActiveUsers": dailyActiveUsers})
}