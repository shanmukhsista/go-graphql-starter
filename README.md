## Introduction

This starter allows you to quickly bootstrap a Golang app that leverages Graphql ( powered by gqlgen framework). This is 
my opinionated way of writing a Graphql application quickly by cloning this template. 

A detailed intro on this starter can be found here. 

[https://medium.com/@shanmukhsista/i-created-a-golang-graphql-project-starter-to-build-high-performance-cloud-native-production-apis-41bb55c05f7](https://medium.com/@shanmukhsista/i-created-a-golang-graphql-project-starter-to-build-high-performance-cloud-native-production-apis-41bb55c05f7)

## Getting started 

To Get started, get started, simply click on the `Use this template` on the top right section above the code listing.

For a detailed tutorial on getting started and renaming your entities, please refer to this article. 


## Pre Requisites for Development

Below are  a list of dependencies that are required to get started with this repository. 

1. Wire Dependency Injection 
   a. https://github.com/google/wire
2. TaskFile https://taskfile.dev/
3. Labstack Echo https://github.com/labstack/echo


## Running Graphql Server Application. 

To run the graphql server, execute the following command from repository root. 

```
task run-graphql-server

task: [run-graphql-server] go build -o graphql-server
task: [run-graphql-server] PORT=7777 ./graphql-server -configpath ./config/dev/config.yaml
{"level":"debug","time":"2022-11-18T21:52:55-05:00","message":"Using config path %!s(**string=0xc0002061c8)"}

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.9.0
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:7777
```

If you see the server started message, visit http://localhost:7777/ in your browser to access graphql ui. 

Try  the following mutation to create a note. 


```graphql
mutation{
  createNewNote(input:{title:"Hello, my new Note!" , 
    content:"Some long note content"}){
    id
    title
    content
  }
}
```

## General Guidelines 

All Methods that are prefixed with `Must*` should generate a panic on failure. 

## Error Handling Rules 

* Every error must be of type `appError`. Use the helpers provided within the apperrors.go file to create new ones.
* Each error must have a unique key associated with it. Translations must be specified within the `error_messages.json` file.


### Creating a new Database Migration

This application uses postgres database as an example. To create a new migration , use the following command and 
define your database migration scripts. 

You can manually apply the migrations or use `golang-migrate` to configure automatic migration config. 
 
```
migrate create -ext sql -dir migrations/ -seq initial_notes_schema
```