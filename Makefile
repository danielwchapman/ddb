build:
	go generate

clean:
	rm -f cover.out

test: ut it

# Replace TEST_TABLE and AWS_PROFILE with your own values
# Integration tests table structure must have 1 GSI and 1 LSI with parameters:
#   Table pk: PK, Table sk: SK
#   GSI: indexName: GSI1, pk: GSI1PK, sk: GSI1SK
#   LSI: indexName: LSI1, pk: PK, sk: LSI1SK
it:
	env TEST_TABLE=IntegrationTest INTEGRATION=on AWS_PROFILE=sandbox go test ./... -shuffle=on -v -coverprofile=cover.out

ut:
	go test -v ./... -shuffle=on -v -coverprofile=cover.out
