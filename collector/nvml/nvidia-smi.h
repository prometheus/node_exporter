#ifndef NVIDIA_SMI_H
#define NVIDIA_SMI_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
   unsigned int clock_mem; // GPU memory clock in Mhz
   unsigned int clock_graphics; // GPU graphics clock in Mhz
   unsigned long long throttle; // throttle reason
   unsigned int perf_state; // performance state; 0: max / 15: min
   unsigned int temperature; // GPU temperature in Celsius degrees
   unsigned int util_gpu; // percentage of time during kernels are executing on the GPU.
   unsigned int util_mem; // percentage of time during memory is being read or written.
   unsigned int mem_usage; // percentage of used memory size
   unsigned int id; // device ID
} GPU_Stats;

int get_gpu_stats(int *pCount, GPU_Stats **ppStats);

#ifdef __cplusplus
}
#endif

#endif // NVIDIA-SMI_H
