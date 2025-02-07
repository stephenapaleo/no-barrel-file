# No Barrel File CLI

![GitHub Release](https://img.shields.io/github/v/release/nergie/no-barrel-file?style=for-the-badge) &nbsp;
![GitHub Release Date](https://img.shields.io/github/release-date/nergie/no-barrel-file?display_date=published_at&style=for-the-badge) &nbsp;
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/nergie/no-barrel-file/build-and-deploy.yml?branch=main&style=for-the-badge&label=build)

---

`no-barrel-file` is a CLI tool for removing barrel file imports in JavaScript/TypeScript projects. It efficiently replaces barrel imports with full paths, counts and displays barrel files in folders, and supports alias paths, gitignore rules, and the handling of nested or circular barrel structures.

![Demo GIF](./assets/demo.gif)

Using `no-barrel-file` to remove barrel files from your project can significantly improve your app’s runtime, build-time and test-time performance.
For more details, reading [Speeding up the JavaScript ecosystem - The barrel file debacle](https://marvinh.dev/blog/speeding-up-javascript-ecosystem-part-7/) is recommended.

---

<!-- TOC start (generated with https://github.com/derlin/bitdowntoc) -->

- [No Barrel File CLI](#no-barrel-file-cli)
  - [Features](#features)
  - [**🔧 Installation**](#-installation)
    - [**Installation With Golang (Recommended)**](#installation-with-golang-recommended)
    - [**Installation With Docker**](#installation-with-docker)
  - [Usage](#usage)
    - [General Command Structure](#general-command-structure)
    - [**🔗 Available Commands**](#-available-commands)
    - [**🌟 Flags**](#-flags)
      - [Global Flags](#global-flags)
      - [Replace Command Flags](#replace-command-flags)
    - [**Count barrel files**](#count-barrel-files)
    - [**Display all barrel files**](#display-all-barrel-files)
    - [**Replace imports of a specific barrel file**](#replace-imports-of-a-specific-barrel-file)
    - [**Replace all barrel file imports**](#replace-all-barrel-file-imports)
  - [**Run the CLI inside a Docker container**](#run-the-cli-inside-a-docker-container)
    - [**Count barrel files**](#count-barrel-files-1)
    - [**Display all barrel files**](#display-all-barrel-files-1)
    - [**Replace imports of a barrel file**](#replace-imports-of-a-barrel-file)
    - [**Replace all barrel file imports**](#replace-all-barrel-file-imports-1)
  - [License](#license)

## <!-- TOC end -->

## Features

- **Display Barrel Files**: List all barrel files in your project for easy inspection.
- **Replace Barrel Imports**: Automatically replace barrel file imports with full paths in your project.
- **Count Barrel Files**: Get the total number of barrel files in your project.
- **Customizable**: Supports root path, alias configurations, gitignore rules, file extensions, files to ignore.
- **Scalable**: Designed to work with large projects with complex barrel file structures (Nested barrel files, circular dependencies).

---

## **🔧 Installation**

### **Installation With Golang (Recommended)**

Ensure you have Go installed (`go1.23+` recommended).

```sh
go install github.com/nergie/no-barrel-file@latest
```

### **Installation With Docker**

Ensure you have Docker installed.

```sh
docker pull nergie42/no-barrel-file:latest
```

---

## Usage

### General Command Structure

```sh
no-barrel-file [command] [flags]
```

### **🔗 Available Commands**

| Command                  | Description                                                        |
| ------------------------ | ------------------------------------------------------------------ |
| `no-barrel-file count`   | Count the number of barrel files in the specified root path.       |
| `no-barrel-file display` | Display all barrel files in the specified root path.               |
| `no-barrel-file replace` | Replace barrel imports with full paths in the specified root path. |

### **🌟 Flags**

#### Global Flags

| Flag                   | Description                                               | Default             |
| ---------------------- | --------------------------------------------------------- | ------------------- |
| `--root-path, -r`      | Root path of the targeted project. **Required**           | None (required)     |
| `--extensions, -e`     | Comma-separated list of file extensions to process.       | `.ts,.js,.tsx,.jsx` |
| `--gitignore-path, -g` | Relative path to `.gitignore` file to apply ignore rules. | `.gitignore`        |
| `--ignore-paths, -i`   | Comma-separated list of directories or files to ignore.   | None                |

#### Replace Command Flags

| Flag                      | Description                                                                                                  | Default |
| ------------------------- | ------------------------------------------------------------------------------------------------------------ | ------- |
| `--alias-config-path, -a` | Relative path to `tsconfig.json` or `jsconfig.json` for alias resolution. **Only JSON files are supported.** | None    |
| `--barrel-path, -b`       | Relative path of a barrel file import to replaced.                                                           | `.`     |
| `--target-path, -t`       | Relative path where imports should be replaced.                                                              | `.`     |
| `--verbose, -v`           | Enable verbose output for detailed logs.                                                                     | None    |

### **Count barrel files**

```sh
no-barrel-file count --root-path .
```

### **Display all barrel files**

```sh
no-barrel-file display --root-path .
```

### **Replace imports of a specific barrel file**

```sh
no-barrel-file replace --root-path . --alias-config-path tsconfig.json --barrel-path src/index.ts
```

### **Replace all barrel file imports**

```sh
no-barrel-file replace --root-path . --alias-config-path tsconfig.json
```

---

## **Run the CLI inside a Docker container**

[Docker registry](https://hub.docker.com/r/nergie42/no-barrel-file)  
:warning: **It is recommended to install using Go instead of Docker, as Docker can increase execution time when mounting directories with many files.**

### **Count barrel files**

```sh
# Navigate to the root directory of your project
cd root-path
docker run --rm \
  -v $(pwd):/app \
  --platform linux/amd64 \
  --user $(id -u):$(id -g) \
  nergie42/no-barrel-file count \
  --root-path /app
```

### **Display all barrel files**

```sh
docker run --rm \
  -v $(pwd):/app \
  --platform linux/amd64 \
  --user $(id -u):$(id -g) \
  nergie42/no-barrel-file display \
  --root-path /app
```

### **Replace imports of a barrel file**

```sh
docker run --rm \
  -v $(pwd):/app \
  --platform linux/amd64 \
  --user $(id -u):$(id -g) \
  nergie42/no-barrel-file replace \
  --root-path /app \
  --alias-config-path tsconfig.json \
  --barrel-path src/index.ts -v
```

### **Replace all barrel file imports**

```sh
docker run --rm \
  -v $(pwd):/app \
  --platform linux/amd64 \
  --user $(id -u):$(id -g) \
  nergie42/no-barrel-file replace \
  --root-path /app \
  --alias-config-path tsconfig.json
```

---

## License

This project is licensed under the [MIT License](LICENSE).
