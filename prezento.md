# Prezento

> Presenting presentations presently with Prezento

A Terminal UI (TUI) tool by:

- Lux (`@k-nox`)
- Michael (`@reesmichael1`)
- Eric (`@ericrallen`)
- Andie (`@bugwhisperer418`)

---

# Code & Commands

Inline code like `echo "Hello, World!"` and inline commands, like `$>pwd`, too!

Commands are code blocks prefixed with `$>`

---

# Getting Started

```$>bash
ping -c 2 8.8.8.8
```

---

# Why?

- Everyone has a terminal
- Markdown renders everywhere and is easy to version control and share
- Small footprint
- Have access to all of the processes available to the terminal
  - **Run Code**
  - **Make System Calls**
  - **Delete your filesystem** - but don't do that

---

# Markdown Rendering

## Headers

Paragraphs and list:

- Like this
- And this

### And More Headers

> And even some sweet quotes!

1. **BOLD**
2. _Italics_
3. ~~StrikeThrough~~
4. [Links](https://github.com/reesmichael1/prezento)

---

# Tables

| Key     | Value |
| ------- | ----- |
| foo     | bar   |
| lorem   | ipsum |
| testing | 1234  |

---

# Syntax Highlighting

```go
switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return p, tea.Quit
		}
  case tea.WindowSizeMsg:
  footerHeight := lipgloss.Height(p.paginator.View())

  if !p.ready {
  	// Since this program is using the full size of the viewport we
  	// need to wait until we've received the window dimensions before
  	// we can initialize the viewport. The initial dimensions come in
  	// quickly, though asynchronously, which is why we wait for them
  	// here.
  	p.viewport = viewport.New(msg.Width, msg.Height-footerHeight)
  	// m.viewport.YPosition = 0
  	p.ready = true
  } else {
  	p.viewport.Width = msg.Width
  	p.viewport.Height = msg.Height - footerHeight
  }
}
```

---

# How?

- [Charm.sh](https://charm.sh/): making the command line glamorous
- [creack/pty](https://github.com/creack/pty): pseudo-terminal

---

# Future Improvements

- Actually run the code interactively
- Watch mode - rebuild presentation on file change
- Configurable styles and syntax highlighting theme

---

# Fin

Thanks!
