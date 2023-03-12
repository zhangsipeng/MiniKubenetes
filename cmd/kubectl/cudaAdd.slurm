#!/bin/bash

#SBATCH --job-name=cublas
#SBATCH --partition=dgx2
#SBATCH --output=%j.out
#SBATCH --error=%j.err
#SBATCH -N 1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=6
#SBATCH --gres=gpu:1

ulimit -s unlimited
ulimit -l unlimited

module load  cuda/10.1.243-gcc-8.3.0
nvcc ./minik8sdata/cudaAdd_.cu -o ./binary/cudaAdd_ -lcublas

./binary/cudaAdd_