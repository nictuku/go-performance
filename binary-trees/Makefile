all: binarytrees.gpp-6.gpp_run

binarytrees.gpp-6.gpp_run: binarytrees.gpp-6.c++.o
	/usr/bin/g++ binarytrees.gpp-6.c++.o -o binarytrees.gpp-6.gpp_run -fopenmp -lboost_system 

binarytrees.gpp-6.c++.o:
	/usr/bin/g++ -c -pipe -O3 -fomit-frame-pointer -march=native  -fopenmp binarytrees.gpp -o binarytrees.gpp-6.c++.o
       
