# What?

This is a simple boiler-plate repository for a REST API based on: Golang Fiber framework, GORM, and MySQL.

# How to compile?

To compile the project, you need to have Go installed in your machine. You can download it from [here](https://golang.org/dl/).

```bash
export PATH=$PATH:/usr/local/go/bin
GOPROXY=direct go mod download
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app.exe .
```

There's already a `Dockerfile` that you can use to build the image.


# How to install?

If you're using Docker, you can use the following command to build the image:

```bash
docker build -t go-fiber-boilerplate .
```

There's also a sample service file that you can use to run the app as a Linux service in your machine. 
Check the `deploy-service` folder for the file details.
