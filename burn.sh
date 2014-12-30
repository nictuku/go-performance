#!/bin/bash

set -ux

# If relevant for your system, run this as root:
# for i in /sys/devices/system/cpu/cpu[0-9]*
# do
# 	    echo performance > $i/cpufreq/scaling_governor
# done

echo -n "# ---------------------" >> report.txt
date >> report.txt

for PROCS in 1 2 3 4 5 6 7 8 9 10; do
	for GOROUTINES in 1 2 3 4 8 16; do
		GOMAXPROCS=${PROCS} go run cpuburn.go ${GOROUTINES} >> report.txt
		# Give it time for all connections to close.
		sleep 5
	done
done
