## Introduction

This starter allows you to quickly bootstrap a Golang app that leverages Graphql ( powered by gqlgen framework). 

## Directory Structure


```
- cmd/
- internal/
- pkg/
 - services/
 - common/
```


## Pre Requisites

1. Wire Dependency Injection 
   a. https://github.com/google/wire
2. TaskFile https://taskfile.dev/
3. Labstack Echo https://github.com/labstack/echo


### Testing Graphql 

localhost:7777


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

### General Guidelines 

All Methods that are prefixed with `Must*` should generate a panic on failure. 

### Error Handling Rules 

* Every error must be of type `AppError`
* Each error must have a unique key associated with it. 


### Creating a new Database Migration

```
migrate create -ext sql -dir migrations/ -seq initial_notes_schema
```