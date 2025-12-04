# 📂 Gorder - The File Organizer

**Tired of cluttered directories? `gorder` is a simple, fast, and effective command-line tool that automatically sorts files into folders based on their extensions.**

[![Go Report Card](https://goreportcard.com/badge/github.com/your-repo/gorder)](https://goreportcard.com/report/github.com/your-repo/gorder)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## 🚀 What it does

`gorder` scans the directory it's run from, identifies the extension of each file (e.g., `.txt`, `.pdf`, `.jpg`), creates a new folder for each extension, and moves the files into their corresponding new homes.

**Before `gorder`:**
```
.
├── document.pdf
├── image.jpg
├── notes.txt
├── presentation.pptx
└── report.docx
```

**After `gorder`:**
```
.
├── docx/
│   └── report.docx
├── jpg/
│   └── image.jpg
├── pdf/
│   └── document.pdf
├── pptx/
│   └── presentation.pptx
└── txt/
    └── notes.txt
```

## 📦 Installation

Ensure you have Go installed and your `GOPATH` is set up correctly.

```sh
go install github.com/your-repo/gorder@latest
```
*(Note: You will need to replace `your-repo` with your actual GitHub repository path once you publish it.)*

## 🛠️ Usage

Simply navigate to the directory you want to organize and run the command:

```sh
gorder
```

The program will automatically handle the rest!

## 🔧 Building from Source

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/your-repo/gorder.git
    cd gorder
    ```
2.  **Build the binary:**
    ```sh
    go build -o gorder main.go
    ```
3.  **Run it:**
    ```sh
    ./gorder
    ```

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
