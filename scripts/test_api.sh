#!/bin/bash

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Testando a API do Apple Music Simulator${NC}\n"

# Teste 1: Busca por "Beatles" (todos os tipos)
echo -e "${GREEN}Teste 1: Busca por 'Beatles' (todos os tipos)${NC}"
curl -s "http://localhost:8080/v1/catalog/us/search?term=Beatles" | jq
echo

# Teste 2: Busca por "Queen" apenas músicas
echo -e "${GREEN}Teste 2: Busca por 'Queen' apenas músicas${NC}"
curl -s "http://localhost:8080/v1/catalog/us/search?term=Queen&types=songs" | jq
echo

# Teste 3: Busca por "Pink Floyd" apenas álbuns
echo -e "${GREEN}Teste 3: Busca por 'Pink Floyd' apenas álbuns${NC}"
curl -s "http://localhost:8080/v1/catalog/us/search?term=Pink%20Floyd&types=albums" | jq
echo

# Teste 4: Busca por "Michael Jackson" apenas artistas
echo -e "${GREEN}Teste 4: Busca por 'Michael Jackson' apenas artistas${NC}"
curl -s "http://localhost:8080/v1/catalog/us/search?term=Michael%20Jackson&types=artists" | jq
echo

# Teste 5: Busca inválida (sem termo)
echo -e "${GREEN}Teste 5: Busca inválida (sem termo)${NC}"
curl -s "http://localhost:8080/v1/catalog/us/search" | jq
echo 