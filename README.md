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

## How does it work?

`compose-check-updates` scans the given directory for Docker Compose files. It then reads the images in the services and checks if there are newer versions available.

If newer versions are found, `compose-check-updates` will suggest the updated image tags. You can then choose to update the Docker Compose files with the new image tags.

## Installation

TBD

## Usage

To check for updates just run the following command:

```bash
compose-check-updates ./your/path/to/docker-compose-files
```