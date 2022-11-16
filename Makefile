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
	@echo "B bin/go-mdb$(EXE) ./cmd/go-mdb"
	@go build -o bin/go-mdb$(EXE) ./cmd/go-mdb
install-go: all-go
	@install -d $(DESTDIR)$(PREFIX)/bin
	@echo I bin/go-mdb$(EXE)
	@cp bin/go-mdb$(EXE) $(DESTDIR)$(PREFIX)/bin
clean-go:
	rm -f bin/go-mdb
## -- AUTO-GO --
