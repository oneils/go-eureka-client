package eureka_test

import (
	"fmt"
	"github.com/oneils/go-eureka-client/pkg/eureka"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const (
	eurekaResponseFile = "../../testdata/eureka-applications.xml"
)

func TestClient_FetchAll(t *testing.T) {

	t.Run("Should fetch all Applications from Eureka", func(t *testing.T) {
		data, err := os.ReadFile(eurekaResponseFile)
		assert.NoError(t, err)

		server := createServer(t, data, http.StatusOK)
		defer server.Close()

		client := http.Client{Timeout: 1 * time.Second}
		c := eureka.NewClient(&client, server.URL+"/eureka/")

		allApps, err := c.FetchAll()

		assert.NoError(t, err)
		assert.Lenf(t, allApps, 8, "Expected 8 application but got %d", len(allApps))
	})

	t.Run("Should return error when can't unmarshal response", func(t *testing.T) {
		expectedError := "Failed to unmarshal Eureka response: XML syntax error on line 1: unexpected EOF"

		data := []byte("<xml>")
		server := createServer(t, data, http.StatusOK)
		defer server.Close()

		client := http.Client{Timeout: 1 * time.Second}
		c := eureka.NewClient(&client, server.URL+"/eureka/")

		_, err := c.FetchAll()

		assert.Error(t, err)
		assert.Equalf(t, expectedError, err.Error(), "Expected error message: %s but got %s", expectedError, err.Error())
	})

	t.Run("Should return error when can't connect to Eureka", func(t *testing.T) {
		expectedError := "Failed to fetch Eureka response: Get \"http://localhost:9999/eureka/\": dial tcp 127.0.0.1:9999: connect: connection refused"

		client := http.Client{Timeout: 1 * time.Second}
		c := eureka.NewClient(&client, "http://localhost:9999/eureka/")

		_, err := c.FetchAll()

		assert.Error(t, err)
		assert.Equalf(t, expectedError, err.Error(), "Expected error message: %s but got %s", expectedError, err.Error())
	})
}

func TestClient_FetchIPAddress(t *testing.T) {

	client := http.Client{Timeout: 1 * time.Second}

	t.Run("should return IP for specified Application Name", func(t *testing.T) {
		expectedIps := []string{"10.200.138.107", "10.200.141.241"}

		data, err := os.ReadFile(eurekaResponseFile)
		if err != nil {
			t.Fatal(err)
		}

		server := createServer(t, data, http.StatusOK)
		defer server.Close()

		c := eureka.NewClient(&client, server.URL+"/eureka/")

		ip, err := c.FetchIPAddress("MY-SERVICE-1")

		assert.NoError(t, err)
		assert.Containsf(t, expectedIps, ip, "Expected IP to be one of %v but got %s", expectedIps, ip)
	})

	t.Run("should return error when fetching Eureka response", func(t *testing.T) {
		expectedError := "Failed to fetch Eureka response. Status code: 500 Internal Server Error"

		server := createServer(t, nil, http.StatusInternalServerError)
		defer server.Close()

		c := eureka.NewClient(&client, server.URL+"/eureka/")

		_, err := c.FetchIPAddress("NOTIFICATION-RETARGETING-SERVICE")

		fmt.Println(err.Error())
		assert.Error(t, err)
		assert.Equalf(t, expectedError, err.Error(), fmt.Sprintf("Expected error message to be '%s' but got '%s'", expectedError, err.Error()))
	})

}

func createServer(t *testing.T, data []byte, statusCode int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eureka/" {
			t.Errorf("Expected path to be /eureka/apps/ but got %s", r.URL.Path)
		}

		w.WriteHeader(statusCode)
		_, err := w.Write(data)
		if err != nil {
			t.Error(err)
		}
	}))
	return server
}
