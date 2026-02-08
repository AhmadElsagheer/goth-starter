#!/bin/bash
set -e

# Default values
MODULES_PATH=""
MODULE_OVERRIDE=""
GOKIT_SOURCE="./go-kit"
PROJECT_ROOT="."

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --modules-path) MODULES_PATH="$2"; shift ;;
        --module) MODULE_OVERRIDE="$2"; shift ;;
        --source) GOKIT_SOURCE="$2"; shift ;;
        --project) PROJECT_ROOT="$2"; shift ;;
        --db-uri) DB_URI="$2"; shift ;;
        --schema-path) SCHEMA_PATH="$2"; shift ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

# Check if project root exists
if [ ! -d "$PROJECT_ROOT" ]; then
    echo "Error: Project directory '$PROJECT_ROOT' not found."
    exit 1
fi

# Check if we are in a Go project root
GO_MOD_PATH="$PROJECT_ROOT/go.mod"
if [ ! -f "$GO_MOD_PATH" ]; then
    echo "Error: go.mod not found at '$GO_MOD_PATH'. Please specify a valid project root with --project or run from the root."
    exit 1
fi

# Detect Target Module
TARGET_MODULE=$(grep "^module" "$GO_MOD_PATH" | awk '{print $2}')
if [ -z "$TARGET_MODULE" ]; then
    echo "Error: Could not determine module name from $GO_MOD_PATH"
    exit 1
fi

echo "Detected target module: $TARGET_MODULE"

# Use override module if provided, otherwise default to target module for backend
BACKEND_IMPORT_PATH="${MODULE_OVERRIDE:-$TARGET_MODULE}"
echo "Using backend import path: $BACKEND_IMPORT_PATH"

# Determine Destination Path
if [ -z "$MODULES_PATH" ]; then
    echo "Usage: ./install-go-kit.sh --modules-path <relative_path> [--module <backend_module_path>] [--source <go-kit-source-path>] [--project <project_path>]"
    echo "Example: ./install-go-kit.sh --modules-path backend/server/modules --project ../my-project"
    exit 1
fi

# Ensure source directory exists
# We specifically need go-kit/auth
SOURCE_AUTH_DIR="$GOKIT_SOURCE/auth"

if [ ! -d "$SOURCE_AUTH_DIR" ]; then
    echo "Error: Source directory '$SOURCE_AUTH_DIR' not found."
    echo "Please ensure the 'go-kit' folder is available in the current directory or specify --source."
    exit 1
fi

# Construct Destination Paths
# DEST_DIR needs to be relative to PROJECT_ROOT for mkdir and cp
# But typically user provides path relative to project root.
FULL_DEST_DIR="$PROJECT_ROOT/$MODULES_PATH/auth"
IMPORT_REL_DIR="$MODULES_PATH/auth"

# Construct Import Path
# Import path is "github.com/my/project/backend/server/modules/auth"
AUTH_IMPORT_PATH="$TARGET_MODULE/$IMPORT_REL_DIR"

echo "Installing go-kit/auth from '$SOURCE_AUTH_DIR' to '$FULL_DEST_DIR'"
echo "Auth import path will be: $AUTH_IMPORT_PATH"

# Create destination directory
mkdir -p "$FULL_DEST_DIR"

# Copy files
echo "Copying files..."
cp -R "$SOURCE_AUTH_DIR/"* "$FULL_DEST_DIR/"

# Replace placeholders
echo "Replacing import paths..."
# Find all go files in the destination
find "$FULL_DEST_DIR" -name "*.go" -type f | while read -r file; do
    # Use | as delimiter for sed to handle slashes in paths
    # MacOS compatible sed -i ''
    sed -i '' "s|{{AUTH_MODULE}}|$AUTH_IMPORT_PATH|g" "$file"
    sed -i '' "s|{{BACKEND_MODULE}}|$BACKEND_IMPORT_PATH|g" "$file"
done

# --- Phase 2: Install go-kit/pkg ---

SOURCE_PKG_DIR="$GOKIT_SOURCE/pkg"
FULL_PKG_DEST_DIR="$PROJECT_ROOT/pkg"

if [ -d "$SOURCE_PKG_DIR" ]; then
    echo "Installing go-kit/pkg from '$SOURCE_PKG_DIR' to '$FULL_PKG_DEST_DIR'"
    mkdir -p "$FULL_PKG_DEST_DIR"
    echo "Copying pkg files..."
    cp -R "$SOURCE_PKG_DIR/"* "$FULL_PKG_DEST_DIR/"

    echo "Replacing import paths in pkg..."
    find "$FULL_PKG_DEST_DIR" -name "*.go" -type f | while read -r file; do
        sed -i '' "s|{{BACKEND_MODULE}}|$BACKEND_IMPORT_PATH|g" "$file"
    done
else
    echo "Warning: go-kit/pkg not found at '$SOURCE_PKG_DIR', skipping."
fi

echo "Replacements complete."

# --- Phase 3: sqddl (Optional) ---
# Check if sqddl is needed (if args provided) and run it BEFORE dependencies
if [ -n "$DB_URI" ] && [ -n "$SCHEMA_PATH" ]; then
    echo "Running sqddl setup..."
    
    # 1. Check DB Connectivity
    echo "Verifying database connectivity for $DB_URI..."
    
    # Create a temporary Go program to check connection
    cat <<EOF > check_db.go
package main
import (
	"fmt"
	"net"
	"net/url"
	"os"
	"time"
)
func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	dbURL := os.Args[1]
	u, err := url.Parse(dbURL)
	if err != nil {
		fmt.Printf("Error checking DB: Invalid URL %v\n", err)
		os.Exit(1)
	}
	host := u.Host
	if host == "" {
        fmt.Println("Error checking DB: Invalid host")
        os.Exit(1)
    }
    // If no port is specified, default to 5432 (postgres default)
    if _, _, err := net.SplitHostPort(host); err != nil {
        host = host + ":5432"
    }
    
	conn, err := net.DialTimeout("tcp", host, 5*time.Second)
	if err != nil {
		fmt.Printf("Error checking DB: Unreachable %s (%v)\n", host, err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Database is reachable.")
}
EOF

    if ! go run check_db.go "$DB_URI"; then
        echo "Database is not reachable. Aborting installation."
        rm check_db.go
        exit 1
    else
        rm check_db.go
        
        # 2. Run sqddl
        echo "Database check passed."
        
        FULL_SCHEMA_PATH="$PROJECT_ROOT/$SCHEMA_PATH"
        mkdir -p "$FULL_SCHEMA_PATH"
        
        (cd "$FULL_SCHEMA_PATH" && {
            echo "Generating tables.go in "$FULL_SCHEMA_PATH"..."
            pwd
            if ! sqddl tables -db "$DB_URI" -file tables.go -pkg sqtypes; then
                echo "Error: sqddl command failed."
                exit 1
            fi
            
            if [ ! -f "tables.go" ]; then
                 echo "Error: tables.go was not created."
                 echo "Hint: Are there tables in the database? Did you run the migrations?"
                 exit 1
            fi
            
            # Check if tables.go is empty (less than 50 bytes, basically just package declaration)
            if [ $(wc -c < "tables.go") -lt 50 ]; then
                 echo "Warning: tables.go seems empty. Did you run migrations?"
            fi
        }) || exit 1
    fi
fi

# Install dependencies
echo "Installing dependencies..."
echo "Running go mod tidy in $PROJECT_ROOT..."

# Run go commands in the project root
(cd "$PROJECT_ROOT" && go mod tidy)

echo "go-kit installed successfully!"
