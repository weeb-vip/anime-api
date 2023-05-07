gql:
	go run github.com/99designs/gqlgen generate

generate: mocks gql


mocks:
	echo "Generating mocks"