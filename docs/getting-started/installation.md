# Installation

## Requirements

- Go 1.21+
- Pandoc (for PDF/DOCX export)
- LuaLaTeX (for PDF with custom fonts)

## Install CLI

```bash
go install github.com/grokify/structured-profile/cmd/sprofile@latest
```

## Install as Library

```bash
go get github.com/grokify/structured-profile
```

## Install Pandoc

=== "macOS"

    ```bash
    brew install pandoc
    brew install --cask mactex  # For lualatex
    ```

=== "Ubuntu/Debian"

    ```bash
    sudo apt-get install pandoc texlive-luatex
    ```

=== "Windows"

    Download from [pandoc.org](https://pandoc.org/installing.html)

## Verify Installation

```bash
sprofile --version
pandoc --version
```
