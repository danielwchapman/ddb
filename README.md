# ddb
ddb is a collection of functions that make AWS DynamoDB easier to work with.

###### Install
```sh
go get github.com/danielwchapman/ddb              
```

###### Unit Testing
```sh
go test ./... -shuffle=on -v
```

###### Example Integration Test command:
```
env TEST_TABLE=Users INTEGRATION=on AWS_PROFILE=sandbox go test ./... -shuffle=on -v -coverprofile=cover.out
```
