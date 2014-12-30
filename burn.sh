#!/bin/bash
set -ux

echo -n "# ---------------------" >> report.txt
date >> report.txt

for PROCS in 1 2 3; do
	for GOROUTINES in 1 2 4 8; do
		GOMAXPROCS=${PROCS} go run cpuburn.go ${GOROUTINES} >> report.txt
	done
done
