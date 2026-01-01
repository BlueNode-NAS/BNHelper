#!/bin/bash
# Local build script for BlueNode Helper
# Usage: ./build.sh [rpm|binary|all|clean]

set -e

BINARY_NAME="bluenode-helper"
VERSION=$(grep 'Version = ' version.go | sed 's/.*"\(.*\)".*/\1/')
BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "dev")
DIST_DIR="dist"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

build_binary() {
    print_info "Building binary..."
    print_info "Version: ${VERSION}, Commit: ${GIT_COMMIT}, Date: ${BUILD_DATE}"
    
    go build -v \
        -ldflags "-X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE} -X main.GitCommit=${GIT_COMMIT}" \
        -o ${BINARY_NAME}
    
    chmod +x ${BINARY_NAME}
    
    print_info "Binary built successfully: ${BINARY_NAME}"
    ./${BINARY_NAME} -version
}

build_rpm() {
    print_info "Building RPM package..."
    
    # Check if rpm tools are installed
    if ! command -v rpmbuild &> /dev/null; then
        print_error "rpmbuild not found. Install with: sudo dnf install rpm-build"
        exit 1
    fi
    
    # Check if binary exists
    if [ ! -f "${BINARY_NAME}" ]; then
        print_error "Binary not found. Run './build.sh binary' first"
        exit 1
    fi
    
    # Setup RPM build environment
    print_info "Setting up RPM build environment..."
    mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
    
    # Copy pre-built binary and service file as sources
    print_info "Copying binary and service file..."
    cp ${BINARY_NAME} ~/rpmbuild/SOURCES/${BINARY_NAME}
    cp ${BINARY_NAME}.service ~/rpmbuild/SOURCES/${BINARY_NAME}.service
    
    # Copy spec file
    cp ${BINARY_NAME}.spec ~/rpmbuild/SPECS/
    
    # Build RPM
    print_info "Running rpmbuild..."
    rpmbuild -ba \
        --define "_topdir $HOME/rpmbuild" \
        --define "version ${VERSION}" \
        ~/rpmbuild/SPECS/${BINARY_NAME}.spec
    
    # Copy RPMs to dist directory
    mkdir -p ${DIST_DIR}
    cp ~/rpmbuild/RPMS/*/*.rpm ${DIST_DIR}/ 2>/dev/null || true
    cp ~/rpmbuild/SRPMS/*.rpm ${DIST_DIR}/ 2>/dev/null || true
    
    print_info "RPM packages built successfully:"
    ls -lh ${DIST_DIR}/*.rpm
}

clean() {
    print_info "Cleaning build artifacts..."
    rm -f ${BINARY_NAME}
    rm -rf ${DIST_DIR}
    rm -rf ~/rpmbuild
    print_info "Clean complete"
}

test_binary() {
    print_info "Running tests..."
    go test -v ./...
}

show_help() {
    cat << EOF
BlueNode Helper Build Script

Usage: $0 [COMMAND]

Commands:
    binary      Build the binary only (default)
    rpm         Build RPM package
    all         Build both binary and RPM package
    test        Run Go tests
    clean       Remove build artifacts
    help        Show this help message

Examples:
    $0              # Build binary
    $0 binary       # Build binary
    $0 rpm          # Build RPM package
    $0 all          # Build everything
    $0 clean        # Clean up

EOF
}

# Main logic
case "${1:-binary}" in
    binary)
        build_binary
        ;;
    rpm)
        build_rpm
        ;;
    all)
        build_binary
        build_rpm
        ;;
    test)
        test_binary
        ;;
    clean)
        clean
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac

print_info "Done!"
