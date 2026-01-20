# ESC (Easy SSH Connection)

ESC is a CLI to manage SSH connections and run interactive sessions or remote commands using a clean, structured configuration.

## Features
- Store SSH connections in a single HCL file.
- Aliases for long connection names.
- Execute commands across multiple hosts (optional parallel).
- List and tree view based on name segments.
- Shell completion for bash/zsh.

## Installation
### From source
```bash
go install github.com/sembraniteam/esc/cmd/esc@latest
```

### Local (this repo)
```bash
go install ./cmd/esc
```

## Quick start
```bash
esc add --name prod --host 10.10.10.1 --user ubuntu
esc connect prod
```

## Configuration
ESC stores its configuration in `~/.esc/config.hcl`. The file is created automatically on first use.

Example:
```hcl
settings {
  timeout_seconds = 30000
}

connection "prod" {
  host = "127.0.0.1"
  user = "ubuntu"
  port = 22
  auth {
    type = "agent"
  }
  tags = ["prod", "api"]
}

connection "staging.web.main" {
  host = "127.0.0.1"
  user = "deploy"
  port = 22
  auth {
    type = "key"
    key_path = "/Users/you/.ssh/id_rsa"
    passphrase_env = "ESC_SSH_PASSPHRASE"
  }
}

connection "db" {
  host = "127.0.0.1"
  user = "postgres"
  port = 22
  auth {
    type = "password"
    password_env = "ESC_DB_PASSWORD"
  }
}

alias "p" {
  target = "prod"
}
```

### Auth type
- `agent`: use SSH agent (`SSH_AUTH_SOCK` must be set).
- `key`: use a private key from `key_path`, optional passphrase from `passphrase_env`.
- `password`: use password from the `password_env` env var.

## Main commands
```bash
esc add --name prod.api.main --host 10.10.10.1 --user ubuntu
esc rm prod.api.main
esc mv prod.api.main prod.api.primary
esc edit prod.api.primary
esc ls
esc ls prod.*
esc tree
esc show prod.api.primary
esc connect prod.api.primary
esc exec prod.* "uptime" --parallel --limit 5
```

### Alias
```bash
esc alias add p prod.api.main
esc alias ls
esc alias rm p
```

### Completion
```bash
esc completion bash
esc completion zsh
```

## Connection naming rules
- Max 3 dot-separated segments, for example: `prod.api.main`.
- Each segment must be lowercase alphanumeric.
- Supported patterns for `ls` and `exec`: `*` or `prefix.*`.
