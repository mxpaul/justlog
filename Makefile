test:
	go test -mod=vendor -gcflags=all=-l ./...

bench:
	go test -mod=vendor -bench=. -benchmem
