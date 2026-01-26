# Hako

Hako is a lightweight, high-performance key-value storage engine accessible over HTTP. Built with Go and the [Fiber](https://gofiber.io/) web framework, Hako provides a simple API for managing multiple isolated databases and key-value pairs.

## Features

-   **HTTP API**: Simple RESTful interface for all operations.
-   **Multiple Databases**: Support for creating and managing distinct logical databases.
-   **Persistence**: Optional snapshot capability to save data to disk using Gob encoding.
-   **Graceful Shutdown**: Ensures data integrity by saving snapshots on shutdown.
-   **CLI**: Built-in CLI for easy management and configuration.

## Getting Started

### Installation

#### GitHub Releases

Download the pre-built binary for your system from the [GitHub Releases](https://github.com/kostya-zero/hako/releases) page.

#### From Source

It is recommended to use the latest version of Go compiler to build.

Clone the repository and build the project:

```bash
git clone https://github.com/yourusername/hako.git
cd hako
go build -o hako
```

### Running the Server

You can start the server using the `run` command:

```bash
./hako run
```

By default, the server listens on `:7000`.

## Configuration

Hako can be configured using a JSON file. Use the `--config` flag to specify the path to your configuration file.

```bash
./hako run --config config.json
```

### Configuration Options

| Key               | Type    | Default               | Description                                      |
| ----------------- | ------- | --------------------- | ------------------------------------------------ |
| `address`         | string  | `:7000`               | The address and port to listen on.               |
| `snapshot_enabled`| boolean | `false`               | Enable or disable data persistence via snapshots.|
| `snapshot_file`   | string  | `hako-snapshot.dat`   | The file path where snapshots will be saved.     |

### Example `config.json`

```json
{
  "address": ":3000",
  "snapshot_enabled": true,
  "snapshot_file": "./data/snapshot.gob"
}
```

## API Reference

### Databases

#### List all databases
```http
GET /db
```

#### Create a new database
```http
POST /db/:database_name
```

#### Delete a database
```http
DELETE /db/:database_name
```

### Key-Value Operations

#### List all keys in a database
```http
GET /db/:database_name/keys
```

#### Get a value
```http
GET /db/:database_name/kv/:key
```
Returns the value as raw text in the response body.

#### Set a value
```http
POST /db/:database_name/kv/:key
```
Send the value in the request body (raw text).

#### Delete a key
```http
DELETE /db/:database_name/kv/:key
```

## Persistence

When `snapshot_enabled` is set to `true`, Hako will automatically save the state of the storage to disk:
1.  **Periodically**: Every 30 seconds.
2.  **On Shutdown**: When the server receives an interrupt signal (SIGINT/SIGTERM).

Snapshots are stored in binary format using Go's `encoding/gob`.

## License

This project is licensed under the MIT License.
See the [LICENSE](LICENSE) file for details.
