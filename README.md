# Powerfactors
This is the Powerfactors technical challenge. 

# Requirements  
![](https://upload.wikimedia.org/wikipedia/commons/2/2d/Go_gopher_favicon.svg) Golang version 1.20+ (Since they try to make every change backwards compatible I don't believe running with a previous version will cause a problem)
![](https://pics.freeicons.io/uploads/icons/png/15889022741579517836-48.png) Docker 

# How to run the program
You have to pass as arguments in the below command the desired address and port that this program will run.
If not, then the default values of `127.0.0.1` and `3000` will take place respectively.
```
$ go run cmd/powerfactors/main.go -address="desired_address" -port="desired_port"
```

There is also a Makefile provided in order to simplify many commands.


Then you need either the command line or a client(e.g. Postman) to make requests. There is a Postman collection under the name `Powerfactors.postman_collection.json` ready for use_
with various http requests that cover many various cases.

![](https://media.giphy.com/media/SwImQhtiNA7io/giphy.gif)
