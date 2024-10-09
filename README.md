# Null Movement Script in Go

This project implements a low-level keyboard control application in Go, inspired by a similar script written in AutoHotkey. It ensures that only one movement key per pair (W/S or A/D) is active at a time, avoiding conflicting inputs. This script can be used in gaming scenarios where simultaneous directional key presses should be avoided.

The project utilizes Windows API functions such as `SetWindowsHookEx` and `SendInput` to capture and manipulate keyboard input at a lower level than typical libraries. The implementation also includes various obfuscation techniques for enhanced binary security.

## Features

- **Low-Level Keyboard Hook**: Uses the Windows API to capture and process keyboard events in real-time.
- **Mutual Exclusion of Movement Keys**: Ensures that only one of the keys from each movement pair (W/S or A/D) is held down at a time.
- **Obfuscation Techniques**: Include obfuscation build methods to minimize the risk of detection.
- **Binary Packing**: Can be further obfuscated using `upx` binary compression.

## Installation

### Prerequisites
- Go 1.16 or higher installed
- Windows operating system (for low-level keyboard control)
- Administrator privileges (required for some low-level hooks)

### Download and Build

1. **Clone the repository**:
    ```bash
    git clone https://github.com/badhex/nullMove.git
    cd nullMove
    ```

2. **Build the project**:
    ```bash
    go build -o null_movement.exe main.go
    ```

   **Optional: Build without the console window**:
   ```
   go build -ldflags="-H windowsgui" -o null_movement.exe main.go
   ```   

   **Optional: Obfuscate the binary using `garble`**:
    * Install `garble`:
      ```bash
      go install mvdan.cc/garble@latest
      ```

      Build using `garble`:
       ```bash
       garble build -o null_movement_obfuscated.exe main.go
       ```

   **Optional: Compress and further obfuscate the binary using `upx`**:
      * Install `upx`:
       ```
       choco install upx
       ```

      * Pack the binary:
       ```
       upx --best --ultra-brute null_movement_obfuscated.exe
       ```

## Usage

Simply run the compiled binary:

```bash
null_movement.exe
