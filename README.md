# ticky
Tiny issue tracking in source code

| todo | in-progress |
| --- | --- |
| | |

- directory-based statuses
- markdown-based issues
  ```markdown
  # tag command @korchasa
   
   1. search for semver in git tags
   2. analyze commit messages to decide that number needs to be changed
   3. create git tag with all messages from commits 
  ```
- `kanban board` in README.md (by pre-commit hook)
- custom issue statuses support
- cli tool with simple commands
  ```bash
    Usage:
      ticky [OPTIONS] <command>
    
    Application Options:
      -v, --verbose     Verbose output
          --statuses=   Comma separated statuses (default: todo,in-progress)
          --issues-dir= Issues directory (default: issues)
    
    Help Options:
      -h, --help        Show this help message
    
    Available commands:
      init    Init ticky
      list    List issues
      my      Show my issues
      new     Create new issue
      readme  Generate README.md
  ```
