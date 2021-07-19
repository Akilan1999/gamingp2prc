# Server Module Implementation 

This section focuses on an in-depth understanding of the server module implementation. To
understand the architecture of the server module refer. The server module can be split
into various sections. Each section will provide information on how a certain feature works.

## Web framework
The web framework used for the server module is called Gin. The reason Gin was chosen is due to
its wide use and strong documentation available on the official github repository. The default
port used is 8088. For version 1.0 of the project ,the server needs to keep port 8088 open to
ensure that other clients and servers can detect it. The possible requests available are GET and
POST for this implementation. The possible responses are either a string or json response or a file.
In the majority of routes a string response refers to an error when calling the following routes.
The following sub topics below will talk about the route implemented:

### /server_info
This route is responsible to get information about the specifications of the
server. The response of this route is in json if the call was successful. 

### /50
This route is responsible for returning a randomly generated 50mb file. This is used to
calculate the download speed from the p2p module.

### /IpTable
This route is a POST request that is responsible to update the server IP table
based on the IP table the client provides. Once the server gets the IP table it checks if the
client is also a server. This is done by calling the url http://<client ip>:8088/server_info. If
the server_info route from the client responds back with computer specifications of the
client. Then the server initially appends the clients IP to the struct. After that the IP table
received from the client is uploaded to the struct. Once this is done the server passes the
struct to the peer to peer module function. The peer to peer module function will return the back with the
new struct with the valid server nodes. The server responds back to the new struct as a
json format. If a string is present in the response then there is probably an error on the
server side.

### /startcontainer
This route takes in a GET request with the number of TCP ports to open and
checks whether the docker container should be hooked to the GPU or not. This route talks
to the docker module implemented as a sub module in the server module. More
information on the docker module in section 5.4.3. This route calls docker the module to
start the container for the client. The docker module returns back a struct. This struct is
returned back to the client as the json response. This struct consists of information such as
docker id, ports numbers open , information regarding SSH and VNC connections to the
docker container created when the client created this request.

### /RemoveContainer
This route takes in a GET request as the container ID. Based on the
container ID provided ,it calls the docker module which deletes the container. If the
deletion is successful it returns back a string which says success.

## Server information/ Specification
This section provides information on how the server specifications are read. There are 2 major
implementations. The first implementation mentions how basic information such as RAM usage,
CPU specification are detected and the second implementation mentions how the GPU drivers are
detected and information is extracted. The client has to assume that the server is using default
docker settings in terms of CPU cycles and other parameters.

### Basic Information 
The file name for these functions is called gopsutil.go. This codebase
uses the library gopsutil. Gopsutil has various packages or modules within the library
which have functions implemented to get system information. The following information is
stored in a struct and the function returns that struct. 

```go
type SysInfo struct {
    Hostname string `bson:hostname`
    Platform string `bson:platform`
    CPU      string `bson:cpu`
    RAM      uint64 `bson:ram`
    Disk     uint64 `bson:disk`
    GPU      *Query  `xml: GpuInfo`
}
```
### GPU Information 
The file name for these functions is called GPU.go. This codebase checks
if the Nvidia driver exists and returns the driver information. To do this a shell
command called nvidia-smi is executed. This shell command is executed with a --xml as flag
to ensure that the output is in the XML format. If there is an output as a xml format, that
means there is an nvidia driver installed, and the function just reads the output and stores it
to the struct and returns the GPU information.

```go

type Query struct {
	DriveVersion string `xml:"driver_version"`
	Gpu  Gpu `xml:"gpu"`
}

type Gpu struct{
	GpuName  string `xml:"product_name"`
	BiosVersion string `xml:"vbios_version"`
	FanSpeed string `xml:"fan_speed"`
	Utilization GpuUtilization `xml:"utilization"`
	Temperature GpuTemperature `xml:"temperature"`
	Clock GpuClock `xml:"clocks"`
}

type GpuUtilization struct {
	GpuUsage string `xml:"gpu_util"`
	MemoryUsage string `xml:"memory_util"`
}

type GpuTemperature struct {
	GpuTemp string `xml:"gpu_temp"`
}

type GpuClock struct {
	GpuClock string `xml:"graphics_clock"`
	GpuMemClock string `xml:"mem_clock"`
}
```

## Docker Module 
This section provides information on how the server module interacts with the docker containers.
The server calls 2 routes which either creates or removes the docker container. Docker has a huge
advantage because it takes less than 20 seconds to spin up a new container once it’s built and
executed at least once. For docker operations a separate module/package has been created. The
following subtopics will provide more information on how this package works.

### Docker Api 
For this the api has been taken from the official docker repository. To be more
specific it is the client module in the official docker repository. Docker was built using Go.
During this project Docker functions could be directly called from the docker repository.
The Docker api initially ensures that it can detect the docker environment variables. Once
detected, it can execute various functions from the docker client module. The reason the
docker api was selected was to detect and handle errors better.

### Docker Image
The docker image used to spin up the containers is called
ConSol/docker-headless-vnc-container. The following container was modified to open
SSH ports for an SSH connection. The following docker image runs ubuntu 16. The reason
this image was chosen as a default is because if the client wants to access the container in
the form of a desktop environment. This image would allow the client to do so from just a
browser. 

### Build container 
This function pulls the docker image locally and builds the image. Initially
there is a timeout function to ensure that building the image does not take too long to
build. The next phase would be based on the path to get the DockerFile. The tag name of
the container is set as p2p-ubuntu as default. Once the following is set then the docker
build command is executed.

### Run container 
After building the container it needs to be executed for the user to access
the container and do certain operations. The docker package/module has a function to do
this. The function takes in the docker environment as a parameter and also the docker
struct. The docker struct has information such as the TCP ports which are supposed to be
open and whether the docker container should have the GPU hooked to it or not. Based on
the appropriate information provided ,the docker image gets started. The Image gets
started by interacting with the docker client modules. When hooking the GPU the docker
run command is called from the shell. This is because the docker Api does not support the
GPU module yet. When the container is executed for the first time it takes
more than 10 minutes to build. From the second time onwards it takes only 10 seconds to
run.

### Stop and remove container 
This implementation here ensures that the docker is stopped, and the container is removed. This is to ensure 
it does not utilize server resources when it is not being used, or the task that is intended to be executed is complete. 
To run this function all that is needed is the docker container ID. If the function is successful it returns
a string that says success.

### Ports json file 
This file will help map internal ports inside a container to external ports inside a container. A common example 
would be the SSH port which is port 22 inside the docker container and is mapped to random TCP port outside container 
so that any external machines can directly connect into the container. The below representation mentions of where 
the ports.json file is located and also the format of that file. 
```
|_ <Container name>
        |_ Dockerfile
        |_ description.txt 
        |_ ports.json  // The ports file 
```
Format of the ports.json file 
```
{
  "Port": [
    {
      "PortName": "<Port name>",
      "InternalPort": <internal port>,
      "Type": "<tcp/udp>",
      "Description": "<description about the port>"
    }, ... n
  ]
}
```





