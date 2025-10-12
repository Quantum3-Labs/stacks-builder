#!/bin/bash

# Clarity Calculator Deployment Script
# This script helps deploy the calculator contract to different networks

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Clarity Calculator Deployment Script${NC}"
echo "=================================="

# Check if Clarinet is installed
if ! command -v clarinet &> /dev/null; then
    echo -e "${RED}Error: Clarinet is not installed. Please install Clarinet first.${NC}"
    echo "Visit: https://github.com/hirosystems/clarinet"
    exit 1
fi

# Function to deploy to testnet
deploy_testnet() {
    echo -e "${YELLOW}Deploying to Stacks Testnet...${NC}"
    clarinet deploy --testnet
    echo -e "${GREEN}Deployment to testnet completed!${NC}"
}

# Function to deploy to mainnet
deploy_mainnet() {
    echo -e "${YELLOW}Deploying to Stacks Mainnet...${NC}"
    echo -e "${RED}Warning: This will deploy to mainnet. Make sure you have tested thoroughly!${NC}"
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        clarinet deploy --mainnet
        echo -e "${GREEN}Deployment to mainnet completed!${NC}"
    else
        echo -e "${YELLOW}Deployment cancelled.${NC}"
    fi
}

# Function to run tests
run_tests() {
    echo -e "${YELLOW}Running tests...${NC}"
    clarinet test
    echo -e "${GREEN}Tests completed!${NC}"
}

# Function to check contract
check_contract() {
    echo -e "${YELLOW}Checking contract syntax...${NC}"
    clarinet check
    echo -e "${GREEN}Contract check completed!${NC}"
}

# Main menu
case "${1:-}" in
    "testnet")
        check_contract
        run_tests
        deploy_testnet
        ;;
    "mainnet")
        check_contract
        run_tests
        deploy_mainnet
        ;;
    "test")
        run_tests
        ;;
    "check")
        check_contract
        ;;
    "console")
        echo -e "${YELLOW}Starting Clarinet console...${NC}"
        clarinet console
        ;;
    *)
        echo "Usage: $0 {testnet|mainnet|test|check|console}"
        echo ""
        echo "Commands:"
        echo "  testnet  - Deploy to Stacks testnet (includes check and test)"
        echo "  mainnet  - Deploy to Stacks mainnet (includes check and test)"
        echo "  test     - Run the test suite"
        echo "  check    - Check contract syntax"
        echo "  console  - Start Clarinet console for development"
        echo ""
        echo "Examples:"
        echo "  $0 testnet    # Deploy to testnet"
        echo "  $0 test       # Run tests only"
        echo "  $0 console    # Start development console"
        exit 1
        ;;
esac
