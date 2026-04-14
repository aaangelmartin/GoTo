# VHS tape for goto demo — `brew install vhs && vhs demo/demo.tape`
Output demo/demo.gif

Set Shell "zsh"
Set FontSize 14
Set Width 1100
Set Height 640
Set Theme "Dracula"
Set TypingSpeed 45ms
Set PlaybackSpeed 1.0

Hide
Type "export GOTO_ALIASES=/tmp/demo-aliases.json GOTO_CONFIG=/tmp/demo-config.toml"
Enter
Type "rm -f $GOTO_ALIASES $GOTO_CONFIG && clear"
Enter
Show

Type "# goto — open anything, remember everything"
Sleep 900ms
Enter
Sleep 300ms

Type "goto add gh github.com/aaangelmartin --tag dev --desc 'My GitHub'"
Sleep 400ms
Enter
Sleep 600ms

Type "goto add claude claude.ai --tag ai"
Enter
Sleep 300ms
Type "goto add jira my-co.atlassian.net --tag work"
Enter
Sleep 300ms
Type "goto add hn news.ycombinator.com --tag read"
Enter
Sleep 500ms

Type "goto ls"
Enter
Sleep 1500ms

Type "goto --dry-run gh"
Enter
Sleep 900ms

Type "goto --dry-run githu   # fuzzy"
Enter
Sleep 1200ms

Type "goto"
Enter
Sleep 3500ms

# TUI session — user would hit /, a, e, d, ?, q
# We just exit for the demo
Type "q"

Sleep 800ms
