#!/bin/bash

set -ux

# If relevant for your system, run this as root:
# for i in /sys/devices/system/cpu/cpu[0-9]*
# do
# 	    echo performance > $i/cpufreq/scaling_governor
# done

go build || exit 1

echo -n "# ---------------------" >> report.txt
go version >> report.txt
echo -n "# ----" >> report.txt
date >> report.txt

for PROCS in $(seq 1 2); do
	for GOROUTINES in 1 2 3 4 8 16; do
		GOMAXPROCS=${PROCS} ./go-performance ${GOROUTINES} >> report.txt
		# Give it time for all connections to close.
		sleep 2
	done
done
