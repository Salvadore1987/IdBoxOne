# IdBoxOne - ICAO Passport Reader

IdBoxOne is a Go-based application that implements an ICAO (International Civil Aviation Organization) passport reader system. This application provides secure document verification capabilities by reading and processing electronic passport data according to ICAO standards.

## Features

- **MRZ Processing**: Reads and validates Machine Readable Zone (MRZ) data from passports
- **ICAO Compliance**: Implements ICAO Doc 9303 standards for electronic passport reading
- **BAC Authentication**: Establishes Basic Access Control for secure communication with passport chips
- **Document Data Extraction**: Retrieves personal information from Data Group 1 (DG1)
- **Photo Extraction**: Reads and encodes passport photos from Data Group 2 (DG2)
- **SIM Card Support**: Reads ICCID from SIM cards using contactless readers
- **WebSocket API**: Provides real-time communication via WebSocket connection
- **Multi-Reader Support**: Works with ELYCTIS card readers (both CL and CNT models)

## Architecture

The application consists of several key components:

- **Main Server** (`IdBoxOne.go`): HTTP server with WebSocket endpoint
- **MRTD Library** (`lib/mrtd.go`): Core passport reading functionality
- **MRZ Parser** (`lib/mrz.go`): Machine Readable Zone validation and parsing
- **Personal Info Handler** (`lib/personal_info.go`): Data structure for passport information
- **Response Handler** (`lib/response.go`): API response formatting
- **Scanner Interface** (`lib/scanner.go`): Document scanner communication

## Technical Stack

- **Language**: Go
- **External Libraries**: 
  - Smart card communication (`github.com/ebfe/scard`)
  - WebSocket support (`github.com/gorilla/websocket`)
  - Logging (`github.com/astaxie/beego/logs`)
  - Configuration management (`github.com/astaxie/beego/config`)
- **Hardware**: ELYCTIS card readers, document scanners
- **DLL Integration**: ElyMRTD.dll for specialized passport operations

## Configuration

The application uses `conf/app.ini` for configuration settings, including server port configuration (default: 4000).

## Security Features

- Implements BAC (Basic Access Control) for secure passport chip communication
- Validates MRZ checksums to ensure data integrity
- Secure handling of biometric data and personal information
- Base64 encoding for photo data transmission

## Use Cases

This system is designed for:
- Border control and immigration systems
- Identity verification services
- Government document processing
- Secure access control systems

## Author

**Sagitov Eldar**
- GitHub: [https://github.com/Salvadore1987](https://github.com/Salvadore1987)