SHELL=/bin/bash

EXE = resync

all: $(EXE)

resync:
	@echo "building $@ ..."
	$(MAKE) -s -f make.inc s=static

clean:
	rm -f $(EXE)