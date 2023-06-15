# Powerfactors
This is the Powerfactors technical challenge. 

# Requirements![](https://upload.wikimedia.org/wikipedia/commons/2/2d/Go_gopher_favicon.svg)  
* Golang version 1.20+ (Since they try to make every change backwards compatible I don't believe running with a previous version will cause a problem)
* Docker 

# How to run the program
You have to pass as arguments in the below command the desired address and port that this program will run.
If not, then the default values of `0.0.0.0` (instead of localhost we use this address to bypass the network isolation of docker) and `3000` will take place respectively.
```
$ go run cmd/powerfactors/main.go -address="desired_address" -port="desired_port"
```

There is also a `Makefile` provided in order to simplify many commands.


Then you need either the command line or a client(e.g. Postman) to make requests. Also there is a Postman collection provided under the name `Powerfactors.postman_collection.json` ready for use 
with various http requests that cover many cases.

![](https://media.giphy.com/media/SwImQhtiNA7io/giphy.gif)
