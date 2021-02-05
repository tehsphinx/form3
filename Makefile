list:
	make -qp | awk -F':' '/^[a-zA-Z0-9][^$$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)print A[i]}' | sort -u

test:
	docker-compose up

lint:
	golangci-lint run

depgraph:
	godepgraph -s github.com/tehsphinx/form3 | dot -Tpng -o depgraph.png

doc:
	@echo "Open this link in browser: http://localhost:6060/pkg/github.com/tehsphinx/form3/"; godoc -http=:6060

