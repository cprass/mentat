# Mentat

A spaced repetition learning system implemented in Go using the FSRS (Free Spaced Repetition Scheduler) algorithm.

## Quick Start

```bash
# Build the application
go build -o mentat

# Run interactive TUI
./mentat run

# Specify a custom vault directory
./mentat run --vault ~/my-flashcards

# Review cards in CLI mode
./mentat review

# Show statistics
./mentat stats
```

## Card Format

Mentat uses Markdown files (`.md`) with special single-character markers to define flashcards. Each Markdown file becomes a single deck.

### Markers

- **`q`** - Question marker. Text between `q` and `a` becomes the question.
- **`a`** - Answer marker. Text between `a` and the next `q` (or end of file) becomes the answer.
- **`c`** - Cloze deletion marker. Everything between double braces `((...))` becomes a cloze text.
- **`n`** - Note/comment marker. Text after `n` is ignored by the parser.

There are two modes, multi-line and single-line. How you declare the marker defines which mode is used.

#### 1. Multi-Line content markers

If the marker is on a single line without any other content it creates a multi-line content block where everything from the line after the marker to the line before the next marker or the end of the file is included.

```markdown
q

## What is the capitol of France?

Hint:

1. It has a famous large metal tower.
2. It also has a lot of museums.
3. You can probably buy a baguette there.

a

## And the answer is ...

*Paris*

> Paris n’est pas une ville, c’est un monde.

```

#### 2. Single-Line content markers

If a marker at the start of a line is followed by a space and some text it becomes a single line of content. This can be useful to create a lot more condensed decks and to write decks more quickly.

```markdown
c ((Fear)) is the ((mind))-killer.
q What is spice?
a A psychedelic drug found on the planet of Arrakis.
```

### Markdown formatting

You can use standard Markdown to control the layout and style of your content. Currently the app only works in the terminal, so the content is rendered as plain text and there is very little use for it, but I plan to add a browser-based review mode which will be able to parse and render Markdown.

## Vault Structure

By default, Mentat looks for markdown files in `~/mentat`. You can organize your decks using subdirectories:

```
~/mentat/
├── geography/
│   ├── europe.md
│   └── asia.md
├── math/
│   └── arithmetic.md
└── history/
    └── world-war-2.md
```

Each `.md` file becomes a separate deck, and the filename (without extension) becomes the deck name.

## Review System

Mentat uses the FSRS (Free Spaced Repetition Scheduler) algorithm to optimize review timing. When reviewing cards, you rate your recall:

- **1 (Again)** - You forgot the answer
- **2 (Hard)** - You recalled with difficulty
- **3 (Good)** - You recalled correctly
- **4 (Easy)** - You recalled effortlessly

Based on your ratings, the algorithm schedules the next review to maximize long-term retention.

## Data Persistence

- **Review history**: Stored in `{vault}/reviews.log`
- **Source files**: Markdown files are never modified by the app
- **Card IDs**: Generated from question/cloze content (SHA-256 hash)

If you change a question's text, it creates a new card with a new review history.

### Automatic Git syncing (experimental)

There is automatic Git syncing and it can be enabled in the config settings. On startup the app will pull from the vault repo and on shutdown it will push a new commit if there are changes to the vault. Use with caution!

## Interactive TUI

The `mentat run` command launches an interactive terminal interface:

### Stats Screen
- View all decks and their statistics
- Navigate with `j`/`k` (vim-style)
- Press `r` to review selected deck
- Press `a` to add a new deck (creates `.md` file with frontmatter)
- Press `q` to quit
- Footer displays available keybindings

### Review Screen
- Cards show question first (front)
- Press `space` to reveal the answer (back)
- Rate your recall with `1`, `2`, `3`, or `4`
- Press `q` to return to stats

## Building from Source

```bash
# Clone the repository
git clone <repository-url>
cd mentat

# Build
go build -o mentat

# Run
./mentat run
```

## Dependencies

- Go 1.24.3 or later
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [tview](https://github.com/rivo/tview) - TUI framework
- [go-fsrs](https://github.com/open-spaced-repetition/go-fsrs) - FSRS algorithm

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.