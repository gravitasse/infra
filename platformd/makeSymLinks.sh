#!/bin/bash

if [ "$1" == "onlp" ]; then
	cd pluginManager/onlp/libs
	ln -sf libonlp.so.1 libonlp.so
	cd -
elif [ "$1" == "openBMC" ]; then
	echo "Build Target is OpenBMC"
else
	echo "Build Target is Dummy"
fi
