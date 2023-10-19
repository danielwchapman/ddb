# ddb
ddb is a collection of functions that make AWS DynamoDB easier to work with.

## Testing

Unit Testing
```
go test ./... -shuffle=on -v
```

Example Integration Test command:
```
env TEST_TABLE=Users INTEGRATION=on AWS_PROFILE=sandbox go test ./... -shuffle=on -v -coverprofile=cover.out
```
