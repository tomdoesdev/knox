<!--suppress HtmlDeprecatedAttribute -->
<div align="center">
  <img src="knox.svg" alt="Knox Logo" width="150" height="150">
</div>

# Knox

A secure local development secrets manager designed to prevent accidental commits of sensitive configuration to version control systems.

## Overview

Knox solves the problem of accidentally committing sensitive information (API keys, database credentials, tokens) to git repositories through `.env` files. It provides a project-based approach to managing secrets during development by storing them locally and injecting them into applications at runtime.

## Key Features

- **Prevent Secret Leaks** - Remove secrets from version-controlled `.env` files  
- **Template-Based Injection** - Use Go templates in `.env.template` files with in-memory secret replacement
- **Clean Environment Control** - Only template variables are passed to processes by default
- **Local Storage** - SQLite backend for fast, reliable local storage
- **Development Focused** - Designed for local development, not production
