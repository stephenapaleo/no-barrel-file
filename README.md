# No Barrel File CLI

`no-barrel-file` is a CLI tool for removing barrel file imports in JavaScript/TypeScript projects. It efficiently replaces barrel imports with full paths, counts and displays barrel files in folders, and supports alias paths, gitignore rules, and the handling of nested or circular barrel structures.

![Demo GIF](./assets/demo.gif)

Using `no-barrel-file` to remove barrel files from your project can significantly improve your app’s runtime, build-time and test-time performance.  
For more information, I recommend reading [Speeding up the JavaScript ecosystem - The barrel file debacle](https://marvinh.dev/blog/speeding-up-javascript-ecosystem-part-7/).

---

## Features

- **Display Barrel Files**: List all barrel files in your project for easy inspection.
- **Replace Barrel Imports**: Automatically replace barrel file imports with full paths in your project.
- **Count Barrel Files**: Get the total number of barrel files in your project.
- **Customizable**: Supports root path, alias configurations, gitignore rules, file extensions, files to ignore.
- **Scalable**: Designed to work with large projects with complex barrel file structures (Nested barrel files, circular dependencies).

---

## **⚡ Quick Start**

### **1️⃣ Install with Docker**

```sh
docker pull nergie42/no-barrel-file
```

### **2️⃣ Run the CLI inside a Docker container**

#### **Count barrel files**

```sh
# Navigate to the root directory of your project
cd root-path
docker run --rm -v $(pwd):/app nergie42/no-barrel-file count --root-path /app
```

#### **Display all barrel files**

```sh
docker run --rm -v $(pwd):/app nergie42/no-barrel-file display --root-path /app
```

#### **Replace imports of a barrel file**

```sh
docker run --rm -v $(pwd):/app nergie42/no-barrel-file replace --root-path /app --alias-config-path tsconfig.json --barrel-path src/index.ts -v
```

#### **replace all barrel file imports**

```sh
docker run --rm -v $(pwd):/app nergie42/no-barrel-file replace --root-path /app --alias-config-path tsconfig.json
```

---

## **🔧 Installation (Without Docker)**

If you prefer to install the CLI directly:

### **1️⃣ Install Go**

Ensure you have Go installed (`go1.23+` recommended).

### **2️⃣ Install the CLI**

```bash
go install github.com/nergie/no-barrel-file@latest
```

### **3️⃣ Run the CLI**

```bash
no-barrel-file --help
```

---

## Usage

### General Command Structure

```bash
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

#### replace Command Flags

| Flag                      | Description                                                                                                  | Default |
| ------------------------- | ------------------------------------------------------------------------------------------------------------ | ------- |
| `--alias-config-path, -a` | Relative path to `tsconfig.json` or `jsconfig.json` for alias resolution. **Only JSON files are supported.** | None    |
| `--barrel-path, -b`       | Relative path of a barrel file import to replaced.                                                           | `.`     |
| `--target-path, -t`       | Relative path where imports should be replaced.                                                              | `.`     |
| `--verbose, -v`           | Enable verbose output for detailed logs.                                                                     | None    |

---

## License

This project is licensed under the [MIT License](LICENSE).
