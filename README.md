<p align="center">
  <img src="./logo.png" alt="Beschrapi Logo" width="200">
</p>

<h1 align="center">Compose-Check-Updates</h1>

<p align="center">
  <strong>
Easily update Docker Compose image tags to their latest versions.
  </strong>
</p>

`compose-check-updates` helps you manage and update images in Docker Compose files, similar to how `npm-check-updates` works for a `package.json`. This tool is heavily inspired by `npm-check-updates` and works in a similar way.


## Table of Contents

- [Table of Contents](#table-of-contents)
- [Installation](#installation)
  - [Windows](#windows)
  - [Linux](#linux)
- [Usage](#usage)
- [Flags](#flags)
- [How does it work?](#how-does-it-work)

## Installation

### Windows

1. Download the latest Windows release from the [Releases](https://github.com/Padi2312/compose-check-updates/releases) page.
2. (Optional) Rename the downloaded file to `ccu.exe` for easier usage.
3. Include the path to `ccu.exe` in your PATH environment variable.


### Linux

1. Download the latest Linux release from the [Releases](https://github.com/Padi2312/compose-check-updates/releases) page.
2. (Optional) Rename the downloaded file to `ccu` for easier usage.
3. Make the file executable by running `chmod +x ccu`.
4. Include the path to `ccu` in your PATH environment variable. 
5. Run `ccu` from the terminal to check if the installation was successful.



## Usage

To check for updates in Docker Compose files in the current directory, run:

Check for updates only (default: only checking patch versions):

```bash
ccu
```

Check for updates and update the Docker Compose files:

```bash
ccu -u
```

Check for updates, update the Docker Compose files, and restart the services:

```bash
ccu -u -r
```

You can also control the update behavior by using the flags described below. 

## Flags

> [!IMPORTANT]
> When using `-i` for interactive mode other arguments (except `-d` for directory) will be ignored.


| Flag     | Description                                                  | Default                 |
| -------- | ------------------------------------------------------------ | ----------------------- |
| `-h`     | Show help message                                            | `false`                 |
| `-u`     | Update the Docker Compose files with the new image tags      | `false`                 |
| `-r`     | Restart the services after updating the Docker Compose files | `false`                 |
| `-i`     | Interactively choose which images to update                  | `false`                 |
| `-d`     | Specify the directory to scan for Docker Compose files       | `.` (current directory) |
| `-f`     | Full update mode, checks updates to latest semver version    | `false`                 |
| `-major` | Only suggest major version updates                           | `false`                 |
| `-minor` | Only suggest minor version updates                           | `false`                 |
| `-patch` | Only suggest patch version updates                           | `true`                  |


## How does it work?

`compose-check-updates` scans the given directory for Docker Compose files. It then reads the images in the services and checks if there are newer versions available.

If newer versions are found, `compose-check-updates` will suggest the updated image tags. You can then choose to update the Docker Compose files with the new image tags.

> [!NOTE]
> All subdirectories are scanned recursively for Docker Compose files.