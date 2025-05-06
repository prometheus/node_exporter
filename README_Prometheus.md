# Prometheus Node Exporter: Comprehensive System and Hardware Metrics Collection for *NIX Environments

## Project Overview

Node Exporter is a powerful Prometheus exporter designed to collect and expose comprehensive system and hardware metrics from *NIX kernel-based operating systems. It provides a flexible, pluggable architecture for gathering detailed system-level performance and resource utilization data.

### Key Features

- **Cross-Platform Compatibility**: Supports multiple operating systems including Linux, Darwin, FreeBSD, NetBSD, OpenBSD, and Solaris
- **Extensive Metric Collection**: Provides over 50 collectors for capturing system metrics such as:
  - CPU statistics
  - Memory usage
  - Disk I/O performance
  - Network interface statistics
  - Hardware monitoring
  - System load information
  - File system metrics

### Core Benefits

- **Granular System Monitoring**: Offers deep insights into system resources and performance
- **Flexible Configuration**: Allows enabling/disabling specific collectors based on monitoring needs
- **Prometheus Integration**: Seamlessly exports metrics in a format compatible with Prometheus monitoring
- **Low Overhead**: Lightweight Go implementation with minimal system impact
- **Customizable**: Supports additional metrics through text file collector and runtime configuration

Node Exporter serves as a critical tool for system administrators and DevOps professionals who need comprehensive, real-time visibility into system performance and resource utilization across various *NIX environments.

## Getting Started, Installation, and Setup

### Prerequisites

- A supported Unix-like operating system (*NIX)
- Go compiler (version 1.16 or higher recommended)
- Basic system administration skills

### Quick Start

1. **Download and Install:**
   ```bash
   git clone https://github.com/prometheus/node_exporter.git
   cd node_exporter
   make build
   ```

2. **Run Node Exporter:**
   ```bash
   ./node_exporter
   ```
   By default, Node Exporter listens on port 9100 and exposes system metrics.

### Installation Options

#### Binary Installation

1. Download the latest release from the [GitHub Releases](https://github.com/prometheus/node_exporter/releases) page.
2. Extract the tarball:
   ```bash
   tar xvfz node_exporter-*.tar.gz
   cd node_exporter-*
   ./node_exporter
   ```

#### Docker Deployment

```bash
docker run -d \
  --net="host" \
  --pid="host" \
  -v "/:/host:ro,rslave" \
  quay.io/prometheus/node-exporter:latest \
  --path.rootfs=/host
```

#### System Service Installation

For systemd-based systems, use the provided service file:
```bash
# Copy the systemd service file
sudo cp examples/systemd/node_exporter.service /etc/systemd/system/
sudo systemctl enable node_exporter
sudo systemctl start node_exporter
```

### Configuration

Node Exporter provides numerous configuration options:

- Use `--help` to view all available flags
- Enable/disable specific metrics collectors
- Configure text file directory for custom metrics
- Set up TLS endpoint (experimental)

#### Example Configuration

```bash
./node_exporter \
  --collector.disable-defaults \
  --collector.cpu \
  --collector.meminfo \
  --collector.filesystem \
  --collector.textfile.directory=/path/to/textfile/metrics
```

#### TLS Support (Experimental)

Create a web configuration file (`web-config.yml`) and start Node Exporter with:
```bash
./node_exporter --web.config.file=web-config.yml
```

### Verification

1. Open `http://localhost:9100/metrics` in a web browser
2. Verify metrics are being exposed
3. Configure Prometheus to scrape the Node Exporter's metrics endpoint

### Platform Support

- Primarily supports Linux systems
- Limited support for BSD variants (FreeBSD, NetBSD, OpenBSD)
- macOS and Solaris have partial functionality

**Note:** Certain collectors have platform-specific limitations and may require additional system configuration.

## Usage

Node Exporter is a Prometheus exporter for hardware and OS metrics. It listens on port 9100 by default and exposes system metrics via an HTTP endpoint.

### Basic Usage

To start the Node Exporter with default settings:

```bash
./node_exporter
```

This will start the exporter on `http://localhost:9100/metrics`, exposing a wide range of system metrics.

### Common Command-Line Options

#### Changing the Listen Port

By default, Node Exporter runs on port 9100. To change the port:

```bash
./node_exporter --web.listen-address=":9200"
```

#### Selecting Collectors

Enable or disable specific metric collectors using flags:

```bash
# Disable all default collectors and enable only CPU and memory metrics
./node_exporter --collector.disable-defaults --collector.cpu --collector.meminfo

# Exclude specific collectors
./node_exporter --collector.netdev.device-exclude="docker.*"
```

### Filtering Metrics Collection

When scraping metrics with Prometheus, you can filter which collectors to include or exclude:

```yaml
# Prometheus configuration to collect only CPU and memory metrics
scrape_configs:
  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']
    params:
      collect[]:
        - cpu
        - meminfo
```

### Examples of Metric Collections

1. View all enabled collectors:
```bash
./node_exporter -h
```

2. Run with textfile collector (for custom metrics):
```bash
./node_exporter --collector.textfile.directory=/path/to/custom/metrics
```

3. Enable experimental or high-cardinality collectors:
```bash
./node_exporter --collector.perf --collector.pressure
```

### Accessing Metrics

After starting Node Exporter, you can view the raw metrics by visiting `http://localhost:9100/metrics` in a web browser or using a tool like `curl`:

```bash
curl http://localhost:9100/metrics
```

### Security Considerations

- By default, Node Exporter runs on all network interfaces
- Use firewall rules to restrict access
- For secure deployments, consider using the experimental TLS support with a web configuration file

## Command Reference

### Basic Command Structure

The Node Exporter is run with various optional flags to customize its behavior:

```
node_exporter [flags]
```

### Key Flags

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `--web.listen-address` | Address to listen on for web interface and telemetry | `:9100` | `--web.listen-address=:9200` |
| `--web.telemetry-path` | Path under which to expose metrics | `/metrics` | `--web.telemetry-path=/node-metrics` |
| `--web.config.file` | Path to web configuration file | `""` | `--web.config.file=web-config.yml` |
| `--path.rootfs` | Path to root filesystem for container monitoring | `/` | `--path.rootfs=/host` |

### Collector Management Flags

#### Enabling/Disabling Collectors

| Flag Pattern | Description | Example |
|--------------|-------------|---------|
| `--collector.<name>` | Enable a specific collector | `--collector.cpu` |
| `--no-collector.<name>` | Disable a specific collector | `--no-collector.diskstats` |
| `--collector.disable-defaults` | Disable all default collectors | |
| `--collector.enable-defaults` | Enable all default collectors | |

#### Collector-Specific Include/Exclude Flags

Some collectors support fine-grained filtering:

| Collector | Include Flag | Exclude Flag | Example |
|-----------|--------------|--------------|---------|
| `arp` | `--collector.arp.device-include` | `--collector.arp.device-exclude` | `--collector.arp.device-include=eth0` |
| `filesystem` | `--collector.filesystem.mount-points-include` | `--collector.filesystem.mount-points-exclude` | `--collector.filesystem.mount-points-exclude=^/(dev\|proc\|sys)($\|/)` |
| `systemd` | `--collector.systemd.unit-include` | `--collector.systemd.unit-exclude` | `--collector.systemd.unit-include=docker.service` |

### Textfile Collector Flags

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `--collector.textfile.directory` | Directory to read text metrics from | `""` | `--collector.textfile.directory=/var/lib/node_exporter/textfile_collector` |

### Perf Collector Flags

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `--collector.perf.cpus` | List of CPUs to collect metrics from | Runtime CPUs | `--collector.perf.cpus=2-6` |
| `--collector.perf.tracepoint` | Tracepoint to collect metrics for | `""` | `--collector.perf.tracepoint="sched:sched_process_exec"` |

### Sysctl Collector Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--collector.sysctl.include` | Sysctl numeric values to expose | `--collector.sysctl.include=vm.user_reserve_kbytes` |
| `--collector.sysctl.include-info` | Sysctl string values to expose | `--collector.sysctl.include-info=kernel.core_pattern` |

### Other Notable Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--help` | Show help and exit | |
| `--version` | Show version information | |

### Example Complete Command

```bash
node_exporter \
  --web.listen-address=:9100 \
  --collector.disable-defaults \
  --collector.cpu \
  --collector.meminfo \
  --collector.filesystem.mount-points-exclude=^/(dev|proc|sys)($|/)
```

**Note:** Always refer to `./node_exporter -h` for the most up-to-date flag information.

## Configuration

The Node Exporter supports various configuration options through command-line flags and a web configuration file.

### Command-Line Configuration

Node Exporter can be configured using numerous command-line flags that control collector behavior, web interface, and other settings. Key configuration options include:

#### Collector Selection
- `--collector.<name>`: Enable a specific collector
- `--no-collector.<name>`: Disable a specific collector
- `--collector.disable-defaults`: Disable all default collectors
- `--collector.enable-defaults`: Enable default collectors

#### Collector-Specific Flags
Some collectors support additional configuration flags to include or exclude specific resources:

- Filesystem collector:
  ```
  --collector.filesystem.mount-points-exclude=^/(dev|proc|sys)($|/)
  --collector.filesystem.fs-types-include=ext4,xfs
  ```

- Network device collector:
  ```
  --collector.netdev.device-include=eth0,enp1s0
  --collector.netdev.device-exclude=docker.*,veth.*
  ```

#### Textfile Collector
The textfile collector requires specifying a directory for metrics files:
```
--collector.textfile.directory=/path/to/node_exporter/textfile_directory
```

### Web Configuration File

An experimental web configuration file supports TLS configuration:

```yaml
# web-config.yml example
tls_server_config:
  cert_file: /path/to/cert.pem
  key_file: /path/to/key.pem
```

Apply the web configuration with:
```
./node_exporter --web.config.file=web-config.yml
```

### Sysctl Collector Configuration

Configure sysctl metric collection:
```
--collector.sysctl.include=vm.user_reserve_kbytes
--collector.sysctl.include-info=kernel.core_pattern
```

### Perf Collector Configuration

Configure CPU and tracepoint metrics collection:
```
--collector.perf --collector.perf.cpus=2-6
--collector.perf.tracepoint="sched:sched_process_exec"
```

Note: Configuration options may vary based on your specific use case and system requirements. Always refer to the `--help` output for the most up-to-date configuration options.

## Project Structure

The Node Exporter project is organized into several key directories and files to support its functionality as a Prometheus exporter for system metrics:

### Main Project Structure
- `collector/`: Contains the core collection logic for various system metrics
  - Individual collector files for different system components (e.g., `cpu_linux.go`, `meminfo_linux.go`)
  - Platform-specific implementations (Linux, BSD, Darwin, etc.)
  - Utility files and test fixtures

- `docs/`: Documentation and additional resources
  - Contains upgrade guides, example configurations, and node monitoring mixins
  - `node-mixin/`: Grafana dashboard and alerting configurations

- `examples/`: Deployment configuration examples
  - Initialization scripts for various init systems
  - Systemd service configurations
  - Platform-specific startup scripts

### Key Source Files
- `node_exporter.go`: Main entry point for the Node Exporter application
- `go.mod` and `go.sum`: Go module dependency management
- `Makefile`: Build and development automation

### Test and Development Resources
- `collector/fixtures/`: Test data and mock system files for different collectors
- `end-to-end-test.sh`: Comprehensive test script
- `test_image.sh`: Container image testing script

### Configuration and Metadata
- `.github/`: GitHub-specific workflows and configuration
- `examples/`: Example configuration files
- `example-rules.yml`: Sample monitoring rules

### Documentation Files
- `README.md`: Project overview and documentation
- `CHANGELOG.md`: Version change history
- `CONTRIBUTING.md`: Guidelines for contributing to the project
- `LICENSE`: Project licensing information

## Technologies Used

### Programming Language
- Go (Golang) v1.23.0

### Core Frameworks and Libraries
- Prometheus Client Libraries
  - `github.com/prometheus/client_golang`
  - `github.com/prometheus/client_model`
  - `github.com/prometheus/common`
  - `github.com/prometheus/exporter-toolkit`
  - `github.com/prometheus/procfs`

### System Interaction Libraries
- Command-line Parsing
  - `github.com/alecthomas/kingpin/v2`
- System Utilities
  - `github.com/coreos/go-systemd/v22`
  - `github.com/godbus/dbus/v5`
  - `github.com/hashicorp/go-envparse`

### Network and Protocol Libraries
- Network Utilities
  - `github.com/mdlayher/netlink`
  - `github.com/mdlayher/ethtool`
  - `github.com/mdlayher/wifi`
  - `github.com/jsimonetti/rtnetlink/v2`

### Monitoring and Performance Libraries
- NTP Support: `github.com/beevik/ntp`
- Performance Utilities:
  - `github.com/hodgesds/perf-utils`
  - `github.com/lufia/iostat`

### Container and Deployment
- Base Container Image: Prometheus Busybox
- Container Registry: Quay.io

### Development and Testing Tools
- Standard Go toolchain
- Go modules for dependency management

## Additional Notes

### Performance and Resource Considerations

When using the Node Exporter, be mindful of the performance implications of different collectors:

- Some collectors are disabled by default due to high cardinality, potential performance overhead, or significant resource demands.
- Before enabling additional collectors in a production environment, carefully test their impact.
- Monitor `scrape_duration_seconds` to ensure collection completes within scrape intervals.
- Use `scrape_samples_post_metric_relabeling` to track changes in metric cardinality.

### Security and Permissions

Certain collectors may require specific system permissions or kernel configurations:

- The Perf collector might need adjusting `kernel.perf_event_paranoid` sysctl parameter.
- Some collectors like `slabinfo` may have restricted file permissions (e.g., `/proc/slabinfo` is typically 0400).
- When running in a containerized environment, additional flags and bind mounts are necessary to access host system metrics.

### Compatibility and Platform Support

- Collector support varies across different operating systems.
- Linux typically has the most comprehensive metric collection.
- Some collectors are platform-specific (e.g., ZFS on FreeBSD, Linux, and Solaris).
- For Windows systems, the [Windows Exporter](https://github.com/prometheus-community/windows_exporter) is recommended.

### Metric Collection Flexibility

- Use `--collector.<name>` to enable specific collectors.
- Disable default collectors with `--no-collector.<name>`.
- Enable only specific collectors using `--collector.disable-defaults --collector.<name>`.

### Advanced Configuration

- Some collectors support include/exclude flags for fine-grained metric collection.
- The Textfile collector allows custom metric injection for machine-specific statistics.
- Utilize sysctl collector to expose system configuration as metrics.

### Experimental Features

- TLS endpoint support is available via web configuration file.
- Carefully review [exporter-toolkit web-configuration](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md) for implementation details.

### Deprecation Notice

Some collectors (NTP, Runit, Supervisord) are deprecated and will be removed in the next major release.

## Contributing

We welcome contributions to the Node Exporter project! Here are some guidelines to help you contribute effectively:

### Ways to Contribute

- Report issues or bugs
- Suggest new features or improvements
- Submit pull requests
- Improve documentation

### Contribution Process

1. Discuss significant changes on the [Prometheus developers mailing list](https://groups.google.com/forum/?fromgroups#!forum/prometheus-developers) before beginning work.
2. For trivial fixes or improvements, you can directly create a pull request.
3. When creating a pull request, address the maintainers (see MAINTAINERS.md) in the description.

### Code Guidelines

- Follow the [Go Code Review Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments)
- Adhere to Go formatting and style best practices
- Sign your work by adding a `Signed-off-by` line to your commit messages:
  ```
  Signed-off-by: Your Name <your.email@example.com>
  ```

### Collector Implementation Rules

- Focus on exposing machine metrics
- Do not transform metrics in hardware-specific ways
- Metrics should be directly exposed from `/proc` or `/sys`
- Do not require root privileges
- Avoid running external commands
- Use the textfile collector for specialized or complex metrics

### Code of Conduct

We are committed to providing a friendly, safe, and welcoming environment. The project follows the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/main/code-of-conduct.md).

## License

This project is licensed under the Apache License, Version 2.0. You can find the full text of the license in the [LICENSE](LICENSE) file.

#### Key Permissions
- Commercial use
- Modification
- Distribution
- Patent use
- Private use

#### Conditions
- License and copyright notice required
- Changes must be documented
- State changes must be indicated

#### Limitations
- No warranty
- No trademark rights
- Liability limitations apply

For complete details, please refer to the full license text at [http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0).