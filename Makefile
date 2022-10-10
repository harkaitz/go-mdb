DESTDIR=
PREFIX=/usr/local
all:
clean:
install:
## -- AUTO-GO --
all:     all-go
install: install-go
clean:   clean-go
all-go:
	@echo "B bin/go-mdb      ./cmd/go-mdb"
	@go build -o bin/go-mdb      ./cmd/go-mdb
install-go: all-go
	@install -d $(DESTDIR)$(PREFIX)/bin
	@echo I bin/go-mdb
	@cp bin/go-mdb      $(DESTDIR)$(PREFIX)/bin
clean-go:
	rm -f bin/go-mdb
## -- AUTO-GO --
