# go-eureka-client

![go-eureka-client](https://github.com/oneils/go-eureka-client/actions/workflows/go.yml/badge.svg)

Eureka client for Go projects for fetching all registered application in Eureka or fetching any available instance IP by
application's name.

## How to use


## Import the package

```bash
     go get github.com/oneils/go-eureka-client
```

Create Eureka client by specifying http client and Eureka server's URL.

```go
   client := http.Client{Timeout: 1 * time.Second}
    c := eureka.NewClient(&client, server.URL+"/eureka/")
```

## Import the package

```go
    import (
		"github.com/oneils/go-eureka-client/pkg/eureka"
)
```

## Fetch all Eureka applications

```go
   apps, err := c.GetApplications()
   if err != nil {
       log.Fatal(err)
   }
   for _, app := range apps {
       fmt.Println(app.Name)
   }
```

## Fetch Application IP by Application name

```go
   app, err := c.GetApplicationByName("my-app")
   if err != nil {
       log.Fatal(err)
   }
   fmt.Println(app.IPAddr)
```
