# drm collector

The drm collector exposes metrics about drm.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_drm_card_info | Card information | card, memory_vendor, power_performance_level, unique_id, vendor |
| node_drm_gpu_busy_percent | How busy the GPU is as a percentage. | card |
| node_drm_memory_gtt_size_bytes | The size of the graphics translation table (GTT) block in bytes. | card |
| node_drm_memory_gtt_used_bytes | The used amount of the graphics translation table (GTT) block in bytes. | card |
| node_drm_memory_vis_vram_size_bytes | The size of visible VRAM in bytes. | card |
| node_drm_memory_vis_vram_used_bytes | The used amount of visible VRAM in bytes. | card |
| node_drm_memory_vram_size_bytes | The size of VRAM in bytes. | card |
| node_drm_memory_vram_used_bytes | The used amount of VRAM in bytes. | card |
