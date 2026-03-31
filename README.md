# bees

<div align="center">
  <img src="./.github/assets/bees.png" alt="Bees Mascot" width="400" />
</div>

## Problem

[Beads](https://github.com/steveyegge/beads) is a powerful alternative to a collection of .md files that grew to serve multi-agent orchestration. For a developer that likes to keep their hands on the wheel while pairing with an agentic navigator, ~80% of it is dead weight.

## Solution

Bees follows Beads as an alternative to a sea of .md files and drops everything else.

## Install      

```sh                                                                                
go install github.com/w-h-a/bees/cmd/bees@latest        
```                                                                                                        

## Quick Start

```sh
bees init --prefix PROJ              # create a bees.db in the current repo
bees create "Design auth flow" \
  --type task --priority 2           # create an issue
bees list --status open              # see all open issues
bees update PROJ-xxx --assignee me   # assign and refine
bees ready                           # what should I work on next?
```

## Commands

```text
bees init [--stealth] [--prefix]     bees context
bees create "title" [flags]          bees ready [--sort --limit]
bees show <id>                       bees upcoming [--days --assignee]
bees update <id> [flags]             bees search <query>
bees close <id>                      bees dep add <id> --blocks <id>
bees reopen <id>                     bees dep remove <id> <id>
bees list [--status --type ...]      bees dep graph [<id>]
bees import <file.jsonl>             bees comment <id> "text"
bees export [-o file.jsonl]          bees handoff <id> [--done --remaining ...]
bees delete [--closed-before --yes]  bees config set|get|list
                                     bees version
```

## Architecture

### Flowchart

```mermaid
graph TD
  subgraph CLI ["bees CLI (cobra)"]
    CMD[Command Layer]
  end

  subgraph Service ["Service Layer"]
    SVC[Service]
  end

  subgraph Client ["Client Layer"]
    REPO_IF[Repo Interface]
    IMP_IF[Importer Interface]
    EXP_IF[Exporter Interface]
  end

  subgraph Infra ["Infrastructure"]
    SQLITE[SQLite via modernc.org/sqlite]
    DB[(bees.db)]
    BEADS[Beads JSONL Parser]
    JSONW[JSONL Writer]
    JSONL[(.jsonl file)]
  end

  subgraph Domain ["Domain Layer"]
    ISSUE[Issue]
    DEP[Dependency]
    COMMENT[Comment]
    HANDOFF[Handoff]
  end

  CMD --> SVC
  SVC --> Domain
  SVC --> REPO_IF
  SVC --> IMP_IF
  SVC --> EXP_IF
  REPO_IF -.-> SQLITE
  SQLITE --> DB
  IMP_IF -.-> BEADS
  BEADS --> JSONL
  EXP_IF -.-> JSONW
  JSONW --> JSONL
```

### ER Diagram

```mermaid
erDiagram
  issues {
    text id PK
    text title
    text description
    text status "open|in_progress|approved|rejected|closed"
    text type "task|bug|feature|chore|decision|epic"
    int priority "0-4"
    int estimate_mins
    text parent_id FK "self-ref"
  }

  dependencies {
    text issue_id PK "FK → issues"
    text depends_on_id PK "FK → issues"
  }

  labels {
    text issue_id PK "FK → issues"
    text label PK
  }

  comments {
    int id PK
    text issue_id FK
    text body
  }

  handoffs {
    int id PK
    text issue_id FK
    text done
    text remaining
    text decisions
    text uncertain
  }

  issues ||--o{ dependencies : "blocked by"
  issues ||--o{ labels : has
  issues ||--o{ comments : has
  issues ||--o{ handoffs : has
```
