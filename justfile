#!/usr/bin/env just --justfile
run app="knox" *args="":
  go run ./cmd/{{app}} {{args}}

build:
  go build -o ./bin ./cmd/...

test:
  go test -v ./...

cover:
  go test ./...

bench:
    go test -benchmem -bench=. ./...

lint:
    golangci-lint run ./...

format:
    golangci-lint fmt ./...

# Create a new task markdown file with ticket ID
task brief:
    #!/usr/bin/env bash
    set -euo pipefail

    # Create tasks directory if it doesn't exist
    mkdir -p ./.tasks

    # Get the next ticket ID by finding the highest existing ID
    if [ -f "./.tasks/.ticket_counter" ]; then
        ticket_id=$(cat "./.tasks/.ticket_counter")
    else
        ticket_id=0
    fi

    # Increment ticket ID
    ticket_id=$((ticket_id + 1))

    # Save the new ticket ID
    echo $ticket_id > "./.tasks/.ticket_counter"

    # Get current timestamp in format YYYYMMDD_HHMMSS
    timestamp=$(date +"%Y%m%d_%H%M%S")

    # Transform brief to lowercase and replace spaces with underscores
    task_name=$(echo "{{brief}}" | tr '[:upper:]' '[:lower:]' | sed 's/ /_/g')

    # Create filename with ticket ID
    filename="${timestamp}_TK${ticket_id}_${task_name}.md"

    # Create the task file
    cat > "./.tasks/${filename}" << EOF
    # TK${ticket_id}: {{brief}}

    ## Ticket ID
    TK${ticket_id}

    ## Description
    {{brief}}

    ## Status
    - [ ] To Do

    ## Git Branch
    \`TK${ticket_id}-${task_name}\`

    ## Notes

    ## Created
    $(date)
    EOF

        echo "Created task file: ./.tasks/${filename}"
        echo "Ticket ID: TK${ticket_id}"
        echo "Suggested git branch: TK${ticket_id}-${task_name}"
        echo ""
        echo "To create the git branch, run:"
        echo "git checkout -b TK${ticket_id}-${task_name}"

# Create a git branch for a ticket
branch ticket_id:
    #!/usr/bin/env bash
    set -euo pipefail

    # Find the task file with the given ticket ID
    task_file=$(find ./.tasks -name "*TK{{ticket_id}}_*.md" | head -1)

    if [ -z "$task_file" ]; then
        echo "Error: No task found with ticket ID TK{{ticket_id}}"
        exit 1
    fi

    # Extract the task name from the filename
    basename=$(basename "$task_file" .md)
    branch_name=$(echo "$basename" | sed 's/^[0-9]*_//')

    echo "Creating branch: $branch_name"
    git checkout -b "$branch_name"

# Create task and immediately create git branch
start brief:
    #!/usr/bin/env bash
    set -euo pipefail
    
    # Create tasks directory if it doesn't exist
    mkdir -p ./.tasks
    
    # Get the next ticket ID by finding the highest existing ID
    if [ -f "./.tasks/.ticket_counter" ]; then
        ticket_id=$(cat "./.tasks/.ticket_counter")
    else
        ticket_id=0
    fi
    
    # Increment ticket ID
    ticket_id=$((ticket_id + 1))
    
    # Save the new ticket ID
    echo $ticket_id > "./.tasks/.ticket_counter"
    
    # Get current timestamp in format YYYYMMDD_HHMMSS
    timestamp=$(date +"%Y%m%d_%H%M%S")
    
    # Transform brief to lowercase and replace spaces with underscores
    task_name=$(echo "{{brief}}" | tr '[:upper:]' '[:lower:]' | sed 's/ /_/g')
    
    # Create filename with ticket ID
    filename="${timestamp}_TK${ticket_id}_${task_name}.md"
    
    # Create the task file
    cat > "./.tasks/${filename}" << EOF
    # TK${ticket_id}: {{brief}}

    ## Ticket ID
    TK${ticket_id}

    ## Description
    {{brief}}

    ## Status
    - [ ] To Do

    ## Git Branch
    \`TK${ticket_id}-${task_name}\`

    ## Notes

    ## Created
    $(date)
    EOF
    
    echo "Created task file: ./.tasks/${filename}"
    echo "Ticket ID: TK${ticket_id}"
    echo "Suggested git branch: TK${ticket_id}-${task_name}"
    
    # Create the git branch
    branch_name="TK${ticket_id}-${task_name}"
    echo "Creating branch: $branch_name"
    git checkout -b "$branch_name"