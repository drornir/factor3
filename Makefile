include .bingo/Variables.mk

.PHONY: bingo-install-and-link-all
bingo-install-and-link-all: $(BINGO)
	@$(BINGO) get -l
