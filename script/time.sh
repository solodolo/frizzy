#/usr/bin/env bash

if [ ! -f ./frizzy ]; then
	echo "frizzy executable not found"
	exit 1
fi	
/usr/bin/time -f "real: %es, mem: %Mkb" ./frizzy -c config.json
