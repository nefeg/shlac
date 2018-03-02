## These will be provided to the target
NAME=shlac
VERSION=0.6-0ubuntu2
SOURCE=https://github.com/umbrella-evgeny-nefedkin/shlac.git
PPA=ppa:onm/shlac

BUILD = $(shell date +%s)

## Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"



SHELL = /bin/bash
GOPATH = $(shell pwd)
GOBIN = $(GOPATH)/bin
TARGET1 = $(NAME)d
TARGET2 = $(NAME)
BIN	= $(DESTDIR)/usr/bin
CONF = $(DESTDIR)/etc/shlac


export PATH := $(PATH):/usr/local/go/bin

.PHONY: all install uninstall clean ppa configure ppa_test configure_test

all:
	@echo "##################################"
	@echo "# 	Compile binaries"
	@echo "##################################"
	@echo $(LDFLAGS)

	mkdir -p $(GOBIN)
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(LDFLAGS) $(TARGET1) $(TARGET2)

	@echo " *** Done ***"
	@echo ""


install: all

	mkdir -p $(BIN)

	@echo " *** Install server"
	install $(GOBIN)/$(TARGET1) $(BIN)

	@echo " *** Install client"
	install $(GOBIN)/$(TARGET2) $(BIN)

	@echo " *** Install config"
	mkdir -p $(CONF)
	install ./config.json $(CONF)/config.json


uninstall:
	rm -rf $(BIN)/$(TARGET1)
	rm -rf $(BIN)/$(TARGET2)
	rm -rf $(CONF)/config.json


configure_test:
	@echo "##############################"
	@echo "# Compile TEST build"
	@echo "##############################"

	$(MAKE) clean
	./configure_ppa.sh shlac $(VERSION) 1 https://verefkin@bitbucket.org/verefkin/shlac.git

	cd build/tmp; debuild  -S -us -uc

	@echo " *** Build(TEST) is compiled ***"
	@echo ""


configure:
	@echo "###############################"
	@echo "# Compile build"
	@echo "###############################"

	@echo " ==> Cleaning build directory..."
	$(MAKE) clean

	@echo " ==> Configure build..."
	./configure_ppa.sh shlac $(VERSION) 0 https://verefkin@bitbucket.org/verefkin/shlac.git

	### build package (https://help.launchpad.net/Packaging/PPA/BuildingASourcePackage)
	cd build/tmp; debuild -S -sa

	@echo " *** Build is compiled ***"
	@echo ""


ppa: ppa_test configure

	@echo " ==> Uploading to PPA..."
	dput -d $(PPA) $(shell ls build/*.changes)

	$(MAKE) clean


ppa_test: configure_test
	@echo "###############################"
	@echo "#**** 	TESTING BUILD	 ****#"
	@echo "###############################"

	@echo " ==> Unpacking..."
	cd build/; dpkg-source -x *.dsc

	@echo " ==> Testing..."
	cd build/$(NAME)-$(VERSION); dh_auto_test;

	@echo " ==> Compiling..."
	cd build/$(NAME)-$(VERSION); dh_auto_build -a

	@echo " *** Build(TEST) is OK ***"
	@echo ""
	$(MAKE) clean


clean:
	@rm -rf build/
