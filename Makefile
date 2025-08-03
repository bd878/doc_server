all : users docs
.PHONY: all users docs

users :
	protoc ./users/userspb/*.proto \
		--go_out=./users/userspb/ \
		--go-grpc_out=./users/userspb/ \
		--go_opt=paths=import \
		--go-grpc_opt=paths=import \
		--go_opt=module="github.com/bd878/doc_server/users/userspb" \
		--go-grpc_opt=module="github.com/bd878/doc_server/users/userspb";

docs : 
	protoc ./docs/docspb/*.proto \
		--go_out=./docs/docspb/ \
		--go-grpc_out=./docs/docspb/ \
		--go_opt=paths=import \
		--go-grpc_opt=paths=import \
		--go_opt=module="github.com/bd878/doc_server/docs/docspb" \
		--go-grpc_opt=module="github.com/bd878/doc_server/docs/docspb";