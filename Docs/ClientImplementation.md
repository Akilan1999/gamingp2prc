# Client Module Implementation

## Topics
1. [Updating the IP table](#updating-the-IP-table)
2. [Reading server specifications](#reading-server-specifications)
3. [Client creating and removing container](#Client-creating-and-removing-container)
4. [Tracking Containers](#Tracking-Containers)
5. [Grouping Containers](#Grouping-Containers)

This section focuses in depth on how the client module works. The client module is incharge of communicating with
different servers based on the IP addresses provided to the user. The IP addresses are derived
from peer to peer modules. The objective here is how the client module interacts with peer to peer module 
and server module.

### Updating the IP table
The client module calls the peer to peer module to get the local IP table initially, Based on the
servers IP addresses available it calls the speedtest function from the peer to peer module to
update IP addresses with information such as latencies, download and upload speeds. Once this is
done the client module does a Rest Api call to the server to download its IP Table. Once the hops are 
done it writes the appropriate results to the Local IP table. Once this is done it prints out the results. 
To derive parameters such as current the public IP address the url “http://ip-api.com/json/” was called. 
This url returns json response of the current public IP address. This feature will be used in the future 
to ensure that the user's current IP address will not be used for a speed test. 
Clients IP table is updated to the server using a form of type multipart.

### Reading server specifications
The client module calls the route /server_specs and reads the json response. If the json response
was successful then it just calls the pretty print function which just prints the json output in the
terminal.

### Client creating and removing container
The client module uses the servers Rest apis to create and delete containers. To create a container
the client requires 3 parameters being the server ip address, the number of the ports the user
wants to open and if the user wants it connected to the GPU or not. The 3 parameters are sent as a
GET request to the server and the server responds with a json file which has information such as
the container ID, ports open , SSH username, SSH password, VNC username and VNC password.
At the moment the username and password are hard coded from the server side for both SSH and
VNC.
To remove a container the client module only requires the server IP address and the container ID.
The client prints the response from the server Rest api.

### Tracking Containers 
Clients create docker images in multiple machines. This means if the client (i.e user) has many 
containers created there needs to be a way to track them. To track containers there is a file 
called ```trackcontainers.json``` which tracks all the containers running. The snippet below 
show a sample structure of file ```trackcontainer.json```.

```
{
	"TrackContainer": [
		{
			"ID": "<ID>",
			"Container": {<docker.DockerVM struct>},
			"IpAddress": "<IP Address>"
		}
	]
} 
```
The default path to the container tracker is ```client/trackcontainers/trackcontainers.json```. 

### Grouping Containers 
When starting a set container possibility to be able to group them. 
The benefit this would be that when executing plugins the group ID would be enough to execute 
plugin in a set of containers. This provides the possibility to execute repetitive tasks in containers in 
a single cli command. To store groups there is a file called ```grouptrackcontainer.json``` which tracks all
the groups currently present set by the client. The snippet below
show a sample structure of file ```grouptrackcontainer.json```.

```
{
 "Groups": [
  {
   "ID": "grp<Random UUID>",
   "TrackContainer": [{client.TrackContainers struct}]
  }
 ]
}
```
The default path to the container tracker is ```client/trackcontainers/grouptrackcontainer.json```. 

### Note:
The group id will be auto-generated and will will have it's own prefix  in the start which will mostly be ```grp<uuid>```.  
When a container is removed using the command. ```p2prc --rm <ip address> --id <container id>```. It will be automatically deleted from the groups it exists in. 

