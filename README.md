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
  - [Quick](#quick)
  - [System-wide](#system-wide)
    - [Windows](#windows)
    - [Linux](#linux)
- [Usage](#usage)
- [Flags](#flags)
- [How does it work?](#how-does-it-work)
- [Troubleshooting](#troubleshooting)
- [Image tags with only x.y versions](#image-tags-with-only-xy-versions)
  - [No new versions found but there are newer versions available](#no-new-versions-found-but-there-are-newer-versions-available)

## Installation

### Quick 

1. Download the latest Windows release from the [Releases](https://github.com/Padi2312/compose-check-updates/releases) page.
2. Run the following command to check current directory for docker compose image updates:
```bash
ccu-<YOUR_ARCHITECTURE>
```

Example for Windows:
```ps1
ccu-windows-amd64.exe
```

Example for Linux:
```bash
chmod +x ccu-linux-amd64
./ccu-linux-amd64
```

### System-wide

#### Windows

1. Download the latest Windows release from the [Releases](https://github.com/Padi2312/compose-check-updates/releases) page.
2. Rename the downloaded file to `ccu.exe` for easier usage.
   1. (Optional) Add the file's directory to your PATH environment variable. So you can run `ccu` from any directory.
3. Run `ccu.exe -v` from the command prompt to check if the installation was successful.


#### Linux

1. Download the latest Linux release from the [Releases](https://github.com/Padi2312/compose-check-updates/releases) page.
2. Rename the downloaded file to `ccu` for easier usage.
  1. (Optional) Move the file to `/usr/local/bin` to make it available system-wide or just add it to your PATH.
3. Make the file executable by running `chmod +x ccu`.
4. Include the path to `ccu` in your PATH environment variable. 
5. Run `ccu -v` from the terminal to check if the installation was successful.


## Usage

To check for updates in Docker Compose files in the current directory, run:

Check for updates only (default: only checking patch versions):

```bash
ccu
```

Check for updates and update the Docker Compose files:

>[!NOTE]
>When choosing this option, `ccu` will create backups of the original Docker Compose files with the `.ccu` extension.

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

## Troubleshooting

## Image tags with only x.y versions

Some images only have `x.y` versions and no `x.y.z` versions. 
This can lead to the following scenario:

Alpine has the following tags:
- `3.14`
- `3.14.1`
- `3.14.0`

If you are using `3.14` in your Docker Compose file, `ccu` will suggest an update to `3.14.1`.

But for Postgres with the following tags:
- `13`
- `13.3`
- `13.4`

If you are using `13.2` in your Docker Compose file, `ccu` will not suggest an update to `13.4` because it's no valid semver version.

_(This might be changed in the future with an additional flag)_

### No new versions found but there are newer versions available

On default `ccu` checks for patch versions only. 

---

Example:
- Current image tag: `1.0.0`
- Latest image tag: `1.1.0`

Result: No newer patch versions available

---

`ccu` on default will not suggest an update in this case. 

To check for all newer versions, use the `-f` flag:
  
```bash 
ccu -f
```

This will suggest the latest version `1.1.0` as an update.