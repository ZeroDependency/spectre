package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"zerodependency.co.uk/spectre/pkg/spectre"
)

// SpectreCheckTimeout determines how often to query the Spectre Server for Test updates
const SpectreCheckTimeout = time.Second * 60

var cacheCheckTime *time.Time
var testCache []*spectre.Test

// SpectreTest is a helper function to create a Gin Handler that wraps the call in a check for valid Spectre Test
func SpectreTest(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cacheCheckTime == nil || time.Since(*cacheCheckTime) > SpectreCheckTimeout {
			spectreServer := os.Getenv("SPECTRE_SERVER")
			if spectreServer == "" {
				spectreServer = "http://localhost:18080"
			}

			spectreURL := fmt.Sprintf("%v/api/v1/spectre/tests/%v", spectreServer, service)

			res, err := http.Get(spectreURL)
			if err != nil {
				c.Next()
				return
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				c.Next()
				return
			}

			if body == nil {
				c.Next()
				return
			}

			err = json.Unmarshal(body, &testCache)
			if err != nil {
				c.Next()
				return
			}

			now := time.Now()
			cacheCheckTime = &now
		}

		for _, test := range testCache {
			if test.IsTriggered(c) && test.InvocationCount > 0 {
				test.Invoke()
				if test.Response == nil {
					c.AbortWithStatus(test.ResponseCode)
					return
				}

				c.AbortWithStatusJSON(test.ResponseCode, test.Response)
				return
			}
		}

		c.Next()
	}
}
