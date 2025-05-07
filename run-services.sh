#!/bin/bash

# Colors for better output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=======================================${NC}"
echo -e "${GREEN}Freight Rate Gateway - Services Runner${NC}"
echo -e "${BLUE}=======================================${NC}"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker is not installed! Please install Docker and try again.${NC}"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Docker Compose is not installed! Please install Docker Compose and try again.${NC}"
    exit 1
fi

# Build and start the services
echo -e "${GREEN}Building and starting all transportadora services...${NC}"
docker-compose up --build -d

# Check if services are running
if [ $? -eq 0 ]; then
    echo -e "${GREEN}All services started successfully!${NC}"
    echo -e "${BLUE}-------------------------------------${NC}"
    echo -e "${GREEN}Service endpoints:${NC}"
    echo -e "${BLUE}Transportadora A (EntregasRápidas):${NC} http://localhost:6000/api/v1/entregas-rapidas/cotacao"
    echo -e "${BLUE}Transportadora B (LogiFretes):${NC} http://localhost:6001/api/cotacoes"
    echo -e "${BLUE}Transportadora C (MegaShipping):${NC} http://localhost:6002/shipping/quote"
    echo -e "${BLUE}-------------------------------------${NC}"
    echo -e "${GREEN}To view logs, run:${NC} docker-compose logs -f"
    echo -e "${GREEN}To stop all services, run:${NC} docker-compose down"
else
    echo -e "${RED}Failed to start services. Check the logs for more information.${NC}"
    exit 1
fi 