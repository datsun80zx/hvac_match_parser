# HVAC Match Parser

A Go application that matches HVAC equipment combinations against AHRI (Air-Conditioning, Heating, and Refrigeration Institute) certified system configurations. This tool helps identify which combinations of outdoor units, indoor units, furnaces, and air handlers are certified to work together.

## Overview

The HVAC Match Parser reads equipment lists and AHRI certification data from CSV files, then generates all possible equipment combinations and identifies which ones have AHRI certification. This is useful for HVAC contractors, distributors, and manufacturers who need to verify compatible equipment configurations.

## Features

- **Equipment Normalization**: Standardizes model numbers across different equipment types
- **Multi-Brand Support**: Processes equipment from multiple manufacturers
- **System Type Recognition**: Handles various HVAC system configurations:
  - Central AC with Air Handler
  - Central AC with Furnace
  - Heat Pump with Air Handler
  - Heat Pump with Furnace
- **Wildcard Expansion**: Automatically expands wildcard model numbers in AHRI data
- **Cartesian Product Generation**: Creates all possible valid equipment combinations
- **AHRI Certification Matching**: Identifies certified equipment combinations
- **CSV Output**: Generates a comprehensive report of all certified matches

## Prerequisites

- Go 1.24.5 or later

## Installation

Clone the repository:

```bash
git clone https://github.com/datsun80zx/hvac_match_parser.git
cd hvac_match_parser
```

## Input File Format

### Equipment List CSV

The equipment list CSV should contain the following columns:

- **Brand**: Manufacturer name
- **Furnace**: Furnace model numbers
- **Outdoor Unit (ac)**: Air conditioner model numbers
- **Outdoor Unit (hp)**: Heat pump model numbers
- **Evaporator Coil**: Evaporator coil model numbers
- **Air Handler**: Air handler model numbers

### AHRI Certification CSV

The AHRI certification CSV should contain four columns:

1. AHRI Number
2. Outdoor Unit Model Number
3. Indoor Unit Model Number
4. Furnace Model Number

Wildcard characters (`*`) are supported in model numbers and will be automatically expanded.

## Usage

1. Update the file paths in `main.go` to point to your CSV files:

```go
csvFileEquip := "/path/to/your/equipment.csv"
csvFileAHRI := "/path/to/your/ahri_certifications.csv"
```

2. Run the application:

```bash
go run main.go
```

3. The application will generate a file named `certified_hvac_matches.csv` containing all certified equipment combinations.

## Output Format

The output CSV contains the following columns:

- **AHRI Number**: The certification number
- **Brand**: Equipment manufacturer
- **Orientation**: System orientation (if applicable)
- **Type of System**: The system configuration type
- **Outdoor Unit**: Outdoor unit model number
- **Furnace**: Furnace model number (if applicable)
- **Evaporator Coil**: Evaporator coil model number (if applicable)
- **Air Handler**: Air handler model number (if applicable)

## How It Works

1. **Read Equipment Data**: Parses the equipment list CSV and categorizes equipment by type
2. **Normalize Model Numbers**: Truncates model numbers to standard lengths based on equipment type
3. **Read AHRI Data**: Loads AHRI certification records
4. **Build Lookup Map**: Creates a hashmap with wildcard expansion for fast certification lookups
5. **Generate Combinations**: For each brand and system type, generates all possible equipment combinations
6. **Find Matches**: Checks each combination against the AHRI certification database
7. **Output Results**: Writes all certified matches to a CSV file

## Model Number Normalization

The application normalizes model numbers to ensure consistent matching:

- **Air Handler**: 12 characters
- **Evaporator Coil**: 11 characters (with special handling for certain prefixes)
- **Furnace**: 11 characters
- **Condenser (AC)**: 11 characters
- **Condenser (HP)**: 11 characters

## Wildcard Handling

The application supports wildcards in AHRI certification data:

- **Furnace wildcards** (position 1): Expands to both 'R' and 'D' variations
- **Indoor unit wildcards** (positions 2 and second-to-last): 
  - Position 2 always becomes 'P'
  - Second-to-last position expands to 'A', 'B', 'C', and 'D'

## Project Structure

```
hvac_match_parser/
├── main.go                          # Application entry point
├── go.mod                           # Go module definition
├── internal/
│   ├── csv_parser.go               # String normalization and sorting utilities
│   ├── csv_reader.go               # CSV file reading and writing functions
│   ├── matcher.go                  # Equipment combination and matching logic
│   └── data_structures/
│       ├── types_equipment.go      # Equipment type definitions
│       ├── types_csv.go            # Output CSV structure
│       └── types_matching.go       # Matching-related types
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is available under an open source license.

## Contact

For questions or support, please open an issue on the GitHub repository.
