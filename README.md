# simple-git

A lightweight terminal user interface for Git, built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea).

## Installation

### Homebrew

```bash
brew install Calvinnix/tap/simple-git
```

## Usage

Run from within a Git repository:

```bash
simple-git             # Interactive status view
simple-git --hide-help # Start with help bar hidden
simple-git --help      # Show help
simple-git --version   # Show version
```

### Setting up an alias

For convenience, add an alias to your shell configuration (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
alias g='simple-git'
```

Then reload your shell or run `source ~/.bashrc` (or equivalent).

## Views

simple-git has multiple views you can navigate between:

- **Status View** (default) - Stage/unstage files, commit, push
- **Diff View** - View and stage/unstage individual hunks
- **Branches View** - Switch, create, and delete branches
- **Stashes View** - Apply, pop, and drop stashes
- **Log View** - Browse commit history
