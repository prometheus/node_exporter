// +build !darwin
// +build !debug

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include "nvml.h"
#include "nvidia-smi.h"

#define DUMP_ERROR(status) \
    fprintf(stderr, "NVML error: %s at %s:%d\n", nvmlErrorString(status), __FILE__, __LINE__);

// get_gpu_stats() returns array of GPU stats and number of GPU devices
int get_gpu_stats(int *pCount, GPU_Stats **ppStats) {
    unsigned int count = 0;
    nvmlReturn_t status;
    int ret = -1;
    GPU_Stats *pBuf = NULL;

    if ((pCount == NULL) || (ppStats == NULL)) {
        fprintf(stderr, "get_gpu_stats(): parameter error");
        return -1;
    }

    status = nvmlInit();
    if (NVML_SUCCESS != status) {
        DUMP_ERROR(status);
        return -1;
    }
    status = nvmlDeviceGetCount(&count);
    if (NVML_SUCCESS != status) {
        DUMP_ERROR(status);
        goto Error;
    }
    pBuf = (GPU_Stats *)calloc(count, sizeof(GPU_Stats));
    if (pBuf == NULL) {
        goto Error;
    }
    for (int i = 0; i < count; ++i) {
        int devID = i;
        nvmlDevice_t nvmlDevice;
        GPU_Stats *p = pBuf + i;
        p->id = devID;
        status = nvmlDeviceGetHandleByIndex(devID, &nvmlDevice);
        if (NVML_SUCCESS != status) {
            DUMP_ERROR(status);
            goto Error;
        }
        status = nvmlDeviceGetClockInfo(nvmlDevice, NVML_CLOCK_GRAPHICS,
            &(p->clock_graphics));
        if (NVML_SUCCESS != status) {
            DUMP_ERROR(status);
            goto Error;
        }
        status = nvmlDeviceGetClockInfo(nvmlDevice, NVML_CLOCK_MEM,
         &(p->clock_mem));
        if (NVML_SUCCESS != status) {
            DUMP_ERROR(status);
            goto Error;
        }
        status = nvmlDeviceGetCurrentClocksThrottleReasons(nvmlDevice, &(p->throttle));
        if (NVML_SUCCESS != status) {
            DUMP_ERROR(status);
            goto Error;
        }
        status = nvmlDeviceGetPerformanceState(nvmlDevice, &(p->perf_state));
        if (NVML_SUCCESS != status) {
            DUMP_ERROR(status);
            goto Error;
        }
        status = nvmlDeviceGetTemperature(nvmlDevice, NVML_TEMPERATURE_GPU,
            &(p->temperature));
        if (NVML_SUCCESS != status) {
            DUMP_ERROR(status);
            goto Error;
        }

        nvmlUtilization_t util;
        status = nvmlDeviceGetUtilizationRates(nvmlDevice, &util);
        if (NVML_SUCCESS != status) {
            DUMP_ERROR(status);
            goto Error;
        }
        p->util_gpu = util.gpu;
        p->util_mem = util.memory;

        nvmlMemory_t nvmlMemory;
        status = nvmlDeviceGetMemoryInfo(nvmlDevice, &nvmlMemory);
        if (NVML_SUCCESS != status) {
            DUMP_ERROR(status);
            goto Error;
        }
        p->mem_usage = (unsigned int)((100*nvmlMemory.used) / (nvmlMemory.total*1.0));
    }
    ret = 0;
    *ppStats = pBuf;
    *pCount = count;
Error:
    if (ret != 0 && pBuf != NULL) {
        free(pBuf);
    }
    nvmlShutdown();
    return ret;
}

