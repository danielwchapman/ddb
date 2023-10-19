ut:
	go test -v ./... -shuffle=on -v -coverprofile=cover.out

it:
	env TEST_TABLE=Users INTEGRATION=on AWS_PROFILE=sandbox go test ./... -shuffle=on -v -coverprofile=cover.out
