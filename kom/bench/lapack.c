#include<stdio.h>
#include<stdlib.h>
#include<stdint.h>

int LAPACKE_dgels(int, char, int, int, int, double*, int, double*, int);

int32_t rand_ = 1592653589;

// uniform double, same version as k.
double rnd() {
	int32_t r = rand_;
	r ^= (r << 13);
	r ^= (r >> 17);
	r ^= (r << 5);
	rand_ = r;
	return 0.5 + ((double)r) / 4294967295.0;
}

int main(int args, char **argv) {
	if (args != 2)  exit(1);
	int nrhs = atoi(argv[1]);
	int rows = 10000;
	int cols = 300;
	
	double *A = (double *)malloc(8*rows*cols);
	double *B = (double *)malloc(8*rows*nrhs);
	for(int i=0;i<rows*cols;i++) A[i]=rnd();
	for(int i=0;i<rows*nrhs;i++) B[i]=rnd();
	
	int info = LAPACKE_dgels(102, 'N', rows, cols, nrhs, A, rows, B, rows);
	
	double m = 0.0;
	int  off = 0;
	for(int i=0;i<nrhs;i++){
	 off += rows;
	 for(int j=0;j<cols;j++){
	  double x = B[off+j];
	  if(x > m)  m = x;
	 }
	}
	printf("%lf\n", m);
	 
	return info;
}

