package spectre

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var spectreServer = os.Getenv("SPECTRE_SERVER")

// Test represents a single Spectre Test defintion
type Test struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Service         string       `json:"service"`
	URL             string       `json:"url"`
	InvocationCount int64        `json:"invocationCount"`
	Response        *interface{} `json:"response"`
	ResponseCode    int          `json:"responseCode"`
	Trigger         Trigger      `json:"trigger"`
}

// Trigger represents the conditions that must pass to trigger a Spectre Test.
type Trigger struct {
	Headers    map[string]string      `json:"headers"`
	Parameters map[string]string      `json:"parameters"`
	Query      map[string]string      `json:"query"`
	Body       map[string]interface{} `json:"body"`
}

// HasTrigger determines whether there is actually a valid trigger available.
func (t *Trigger) HasTrigger() bool {
	if len(t.Headers) == 0 &&
		len(t.Parameters) == 0 &&
		len(t.Query) == 0 &&
		len(t.Body) == 0 {
		return false
	}

	return true
}

// IsTriggered returns true if the Gin Context matches the Trigger conditions.
func (t *Test) IsTriggered(c *gin.Context) bool {
	if !t.Trigger.HasTrigger() {
		return false
	}

	if t.URL != c.Request.URL.String() {
		return false
	}

	for k, v := range t.Trigger.Headers {
		if header := c.GetHeader(k); header != "" {
			if header != v {
				return false
			}
		}
	}

	for k, v := range t.Trigger.Parameters {
		if param := c.Param(k); param != "" {
			if param != v {
				return false
			}
		}
	}

	for k, v := range t.Trigger.Query {
		if query := c.Query(k); query != "" {
			if query != v {
				return false
			}
		}
	}

	buf, _ := ioutil.ReadAll(c.Request.Body)
	tmpCopy := ioutil.NopCloser(bytes.NewBuffer(buf))
	bodyCopy := ioutil.NopCloser(bytes.NewBuffer(buf))
	c.Request.Body = bodyCopy // This has to set back immediately in case of failure

	readBuffer := new(bytes.Buffer)
	readBuffer.ReadFrom(tmpCopy)

	var tempBody map[string]interface{}
	err := json.Unmarshal(readBuffer.Bytes(), &tempBody)
	if err != nil {
		return false
	}

	for k, v := range t.Trigger.Body {
		if tempBody[k] != v {
			return false
		}
	}

	return true
}

// Invoke marks the Test as invoked and reduces the available invocation count.
func (t *Test) Invoke() {
	if spectreServer == "" {
		spectreServer = "http://localhost:18080"
	}

	spectreURL := fmt.Sprintf("%v/api/v1/spectre/tests/%v/invoke", spectreServer, t.ID)
	http.Post(spectreURL, "application/json", nil)

	t.InvocationCount--
	if t.InvocationCount < 0 {
		t.InvocationCount = 0
	}
}
