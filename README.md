# Powerfactors ![](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAACYUlEQVR4Ab3XA+xcSQAH4N3adhueGZ/tCw7B2bZt2/YFx9qNatu2bXcx/epXrbLvP8kXPP6GmUkUVa4NhJP4mfmsYgj3UCsRa/FzjmcAWULEBl6mapwBqvM9WbbRi9+YRWA558RZ+8tYTZa2NKQSV7OWLG8QS4D6dCawhHPYd685kwl8S7LcP09yG1vI8BlVIvebMoHA93EEaM0wAlM4nkTFBNDHvEiK7Twd/UG8AdSSUyOjfADNEkpFBajGN2TZyI0k8gT4gWS5an8xKwm0pfZRnq1LPwL9aUWSBFBCgLq0J7CMC8g1S14hTYq+vM/bvMWjnE7lYqbdzWwmw1dULWCmdCNNiCLDXO6kSiEBWjKYwDROpJAua82bDGDCXlPZSGAJF5L3Q4+QYgfPU6nIaVuXpnu15M7IWPoh9/c0NW0JjKZVGQZ0TXoQ6EWNHA+76SECPahZhgAtGUWgPdXyNeHXBFZyB61oulfdErrkSXaQ4qlCxsAFLCWwiWlM2GsAr9GmwIF5AlMJDCusS00V7mIOGUIUabrROs93qvIJGTZzR+GrpEWD03ieL/iSb+hHijRvkMxR+3MiLdmFeqUMoKgkrehHoB91j/Jebf4msIrLKdvm5HMCE2h6lNDXs4EsP1GtnLujT/MEaEofArM5nUTFBHCfR9lOitcj0zbmAGrJsUwiMJI2CaWiAlThAzJs4Z64dsgfEZhM80jtz2RxZBmvH9cB5VmyrOFqKtGQv8mymquI7Yh2Dssjo/w7urGNLL9SPc4zYlVeYgMhIsuQ6AYmzhC1eIBhrGI+f3JKsT/fCQ+Uz9zrqGELAAAAAElFTkSuQmCC)
This is the Powerfactors technical challenge. 

# Requirements ![](https://upload.wikimedia.org/wikipedia/commons/2/2d/Go_gopher_favicon.svg)
Golang version 1.20+ (Since they try to make every change backwards compatible I don't believe running with a previous version will cause a problem)

# How to run the program
You have to pass as arguments in the below command the desired address and port that this program will run.
If not, then the default values of `127.0.0.1` and `3000` will take place respectfully.
` go run cmd/powerfactors/main.go -address="desired_address" -port="desired_port"`

Then you need either the command line or a client(e.g. Postman) to make requests. There is a Postman collection under the name `Powerfactors.postman_collection.json` ready for use
with various http requests that cover many various cases.

![](https://media.giphy.com/media/SwImQhtiNA7io/giphy.gif)
