#!/bin/bash

# Usage: ./init-go.sh username/reponame [destination_path]
# Example: ./init-go.sh octocat/my-service ~/projects

MODULE_INPUT=$1
TARGET_DIR=${2:-.}

# Validate input format
if [[ ! "$MODULE_INPUT" =~ .*/.* ]]; then
    echo "Error: First argument must be 'username/reponame'"
    exit 1
fi

USER=$(echo "$MODULE_INPUT" | cut -d'/' -f1)
REPO=$(echo "$MODULE_INPUT" | cut -d'/' -f2)

if [ -d "$TARGET_DIR" ]; then
    echo "Error: Directory '$TARGET_DIR' already exists."
    exit 1
fi

echo "🚀 Initializing project: $REPO"

# 1. Create directory structure
mkdir -p "$TARGET_DIR/cmd/server"
cd "$TARGET_DIR" || exit

# 2. Initialize Go Module
go mod init "github.com/$USER/$REPO"

# 3. Create boilerplate main.go
cat <<EOF > "cmd/server/main.go"
package main

import "fmt"

func main() {
	fmt.Printf("Service $REPO is running.\n")
}
EOF

# 4. Generate docker-compose.yaml
cat <<EOF > "docker-compose.yaml"
version: "3.5"
name: $REPO
services:
  postgres:
    container_name: $REPO-postgres
    image: postgres:15.2
    command: -c 'max_connections=500' -c 'fsync=off' -c 'synchronous_commit=off' -c 'wal_level=minimal' -c 'max_wal_senders=0'
    networks:
      - $REPO-network
    environment:
      - POSTGRES_USER=$REPO
      - POSTGRES_PASSWORD=$REPO
    ports:
      - 5432:5432
    volumes:
      - postgres_data:/var/lib/postgresql/data

networks:
  $REPO-network:
    driver: bridge
    name: $REPO-network

volumes:
  postgres_data: {}
EOF

# 5. Cleanup
go mod tidy
go fmt ./...

echo "✅ Done! Docker and Go files initialized in $TARGET_DIR"
