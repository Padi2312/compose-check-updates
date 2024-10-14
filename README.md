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
- [Usage](#usage)
- [Flags](#flags)
- [How does it work?](#how-does-it-work)

## Installation

TBD


## Usage

To check for updates in Docker Compose files in the current directory, run:

```bash
compose-check-updates 
```

You can also add some flags to customize the behavior:

```bash
compose-check-updates [-u] [-r] [-i] [-d <directory>]
```

See the [Flags](#flags) section for more information.


## Flags

> [!IMPORTANT]
> When using `-i` for interactive mode other arguments (except `-d` for directory) will be ignored.

- `-h` - Show help message
- `-u` - Update the Docker Compose files with the new image tags
- `-r` - Restart the services after updating the Docker Compose files
- `-i` - Interactively choose which images to update
- `-d` - Specify the directory to scan for Docker Compose files


## How does it work?

`compose-check-updates` scans the given directory for Docker Compose files. It then reads the images in the services and checks if there are newer versions available.

If newer versions are found, `compose-check-updates` will suggest the updated image tags. You can then choose to update the Docker Compose files with the new image tags.

> [!NOTE]
> All subdirectories are scanned recursively for Docker Compose files.