# go-eureka-client

Eureka client for Go projects for fetching all registered application in Eureka or fetching any available instance IP by
application's name.

## How to use


## Import the package

```bash
    go get "github.com/go-eureka/eureka"
```

Create Eureka client by specifying http client and Eureka server's URL.

```go
   client := http.Client{Timeout: 1 * time.Second}
    c := eureka.NewClient(&client, server.URL+"/eureka/")
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
