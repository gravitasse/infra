RM=rm -f
RMFORCE=rm -rf 
MKDIR = mkdir -p
DESTDIR=$(SR_CODE_BASE)/snaproute/src/out/bin
SR_PLATFORMD_CODE_BASE=$(SR_CODE_BASE)/snaproute/src/infra/platformd
PARAMSDIR = $(DESTDIR)/params
GENERATED_IPC=$(SR_CODE_BASE)/generated/src
IPC_GEN_CMD=thrift
SRCS=main.go
IPC_SRCS=rpc/platformd.thrift
COMP_NAME=platformd
GOLDFLAGS=-r /opt/flexswitch/sharedlib
PLATFORM_BUILD_TARGET = onlp

#CGO CFLAGS
CGO_CFLAGS=-I$(SR_PLATFORMD_CODE_BASE)/pluginManager/pluginCommon/
ifeq ($(PLATFORM_BUILD_TARGET), openBMC)
export CGO_CFLAGS+=-I$(SR_PLATFORMD_CODE_BASE)/pluginManager/openBMC/inc
export CGO_LDFLAGS+=-L$(SR_PLATFORMD_CODE_BASE)/pluginManager/openBMC/libs
else ifeq ($(PLATFORM_BUILD_TARGET), onlp)
export CGO_CFLAGS+=-I$(SR_PLATFORMD_CODE_BASE)/pluginManager/onlp/inc
export CGO_LDFLAGS+=-L$(SR_PLATFORMD_CODE_BASE)/pluginManager/onlp/libs -lonlp -lm -Wl,--unresolved-symbols=ignore-in-shared-libs
else 
export CGO_CFLAGS+=-I$(SR_PLATFORMD_CODE_BASE)/pluginManager/onlp/inc
export CGO_LDFLAGS+=-L$(SR_PLATFORMD_CODE_BASE)/pluginManager/onlp/libs -lonlp -lm -Wl,--unresolved-symbols=ignore-in-shared-libs
endif


all:ipc exe 
ipc:
	$(IPC_GEN_CMD) -r --gen go -out $(GENERATED_IPC) $(IPC_SRCS)

exe: $(SRCS)
	bash makeSymLinks.sh $(PLATFORM_BUILD_TARGET)
	go build -o $(DESTDIR)/$(COMP_NAME) -ldflags="$(GOLDFLAGS):$(SR_PLATFORMD_CODE_BASE)/pluginManager/onlp/libs" $(SRCS)

guard:
ifndef SR_CODE_BASE
	$(error SR_CODE_BASE is not set)
endif

install:
	$(MKDIR) $(PARAMSDIR)
	install params/platformd.conf $(PARAMSDIR)/
clean:guard
	$(RM) $(DESTDIR)/$(COMP_NAME) 
	$(RMFORCE) $(GENERATED_IPC)/$(COMP_NAME)
