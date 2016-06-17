COMPS=sysd\
	fMgrd

IPCS=sysd\
	fMgrd

all: ipc exe 

exe: $(COMPS)
	 $(foreach f,$^, make -C $(f) exe;)

ipc: $(IPCS)
	 $(foreach f,$^, make -C $(f) ipc;)

clean: $(COMPS)
         $(foreach f,$^, make -C $(f) clean;)

install: $(COMPS)
	 $(foreach f,$^, make -C $(f) install;)

clean: $(COMPS)
	 $(foreach f,$^, make -C $(f) clean;)
