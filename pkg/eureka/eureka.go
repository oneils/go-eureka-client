package eureka

import (
	"encoding/xml"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// AppsResponse is an EUREKA response and maps to xml application tag
type AppsResponse struct {
	Applications []Application `xml:"application" json:"application"`
}

// Application is an EUREKA application and maps to xml instance tag
type Application struct {
	Name      string     `xml:"name" json:"name"`
	Instances []Instance `xml:"instance" json:"instance"`
}

// Instance is an EUREKA instance and maps to xml instance tag
type Instance struct {
	InstanceID     string   `xml:"instanceId" json:"instanceId"`
	HostName       string   `xml:"hostName" json:"hostName"`
	App            string   `xml:"app" json:"app"`
	IPAddr         string   `xml:"ipAddr" json:"ipAddr"`
	Port           string   `xml:"port" json:"port"`
	HomePageUrl    string   `xml:"homePageUrl" json:"homePageUrl"`
	StatusPageUrl  string   `xml:"statusPageUrl" json:"statusPageUrl"`
	HealthCheckUrl string   `xml:"healthCheckUrl" json:"healthCheckUrl"`
	Status         string   `xml:"status" json:"status"`
	Metadata       Metadata `xml:"metadata" json:"metadata"`
}

// Metadata is an EUREKA metadata and maps to xml metadata tag
type Metadata struct {
	JmxPort        string `xml:"jmx.port"`
	ManagementPort string `xml:"management.port"`
}

type client struct {
	httpClient http.Client
	eurekaURL  string
}

const (
	statusUp = "UP"
)

// NewClient creates a new Eureka client
func NewClient(httpClient *http.Client, eurekaURL string) *client {
	return &client{httpClient: *httpClient, eurekaURL: eurekaURL}
}

// FetchAll returns all applications from Eureka
func (c *client) FetchAll() ([]Application, error) {
	response, err := c.httpClient.Get(c.eurekaURL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch Eureka response")
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Failed to fetch Eureka response. Status code: %s", response.Status)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read Eureka response")
	}

	var appsResponse AppsResponse
	err = xml.Unmarshal(data, &appsResponse)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal Eureka response")
	}

	return appsResponse.Applications, nil
}

// FetchIPAddress returns a random instance IP address
func (c *client) FetchIPAddress(appName string) (string, error) {
	allApps, err := c.FetchAll()
	if err != nil {
		return "", err
	}

	// is used to randomize the instance
	rand.Seed(time.Now().UnixNano())

	for _, app := range allApps {
		if app.Name == appName {
			instancesByStatus := make(map[string][]Instance)
			for _, instance := range app.Instances {
				instancesByStatus[instance.Status] = append(instancesByStatus[instance.Status], instance)
			}

			aliveInstances := instancesByStatus[statusUp]

			minIndex := 0
			maxIndex := len(aliveInstances) - 1
			index := rand.Intn(maxIndex-minIndex+1) + minIndex
			return aliveInstances[index].IPAddr, nil
		}
	}
	return "", nil
}
