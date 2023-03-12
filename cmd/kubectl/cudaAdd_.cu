
#include "cuda_runtime.h"
#include "device_launch_parameters.h"
 
#include <stdio.h>
#include <math.h>
#define Row  8
#define Col 4
 
 
__global__ void addKernel(int **C,  int **A)
{
	int idx = threadIdx.x + blockDim.x * blockIdx.x;
	int idy = threadIdx.y + blockDim.y * blockIdx.y;
	if (idx < Col && idy < Row) {
		C[idy][idx] = A[idy][idx] + 10;
	}
}
 
int main()
{
	int **A = (int **)malloc(sizeof(int*) * Row);
	int **C = (int **)malloc(sizeof(int*) * Row);
	int *dataA = (int *)malloc(sizeof(int) * Row * Col);
	int *dataC = (int *)malloc(sizeof(int) * Row * Col);
	int **d_A;
	int **d_C;
	int *d_dataA;
	int *d_dataC;
    //malloc device memory
	cudaMalloc((void**)&d_A, sizeof(int **) * Row);
	cudaMalloc((void**)&d_C, sizeof(int **) * Row);
	cudaMalloc((void**)&d_dataA, sizeof(int) *Row*Col);
	cudaMalloc((void**)&d_dataC, sizeof(int) *Row*Col);
	//set value
	for (int i = 0; i < Row*Col; i++) {
		dataA[i] = i+1;
	}

	for (int i = 0; i < Row; i++) {
		A[i] = d_dataA + Col * i;
		C[i] = d_dataC + Col * i;
	}
	
	cudaMemcpy(d_A, A, sizeof(int*) * Row, cudaMemcpyHostToDevice);
	cudaMemcpy(d_C, C, sizeof(int*) * Row, cudaMemcpyHostToDevice);
	cudaMemcpy(d_dataA, dataA, sizeof(int) * Row * Col, cudaMemcpyHostToDevice);
	dim3 threadPerBlock(4, 4);
	dim3 blockNumber( (Col + threadPerBlock.x - 1)/ threadPerBlock.x, (Row + threadPerBlock.y - 1) / threadPerBlock.y );
	printf("Block(%d,%d)   Grid(%d,%d).\n", threadPerBlock.x, threadPerBlock.y, blockNumber.x, blockNumber.y);
	addKernel << <blockNumber, threadPerBlock >> > (d_C, d_A);

	cudaMemcpy(dataC, d_dataC, sizeof(int) * Row * Col, cudaMemcpyDeviceToHost);
 
	for (int i = 0; i < Row*Col; i++) {
		if (i%Col == 0) {
			printf("\n");
		}
		printf("%d\t", dataC[i]);
	}
	printf("\n");
    
}
