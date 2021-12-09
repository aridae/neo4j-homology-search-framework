#!/bin/bash

# defaults 
URL="https://1001genomes.org/data/GMI-MPI/releases/v3.1/pseudogenomes/fasta/"
DIRECTORY="/fasta"
EXTENSION=".fasta.gz"
PREFIX="pseudo"

POSITIONAL=()
while [[ $# -gt 0 ]]; do
  key="$1"

  case $key in
    -ext|--extension)
      EXTENSION="$2"
      shift # past argument
      shift # past value
      ;;
    -prefix|--prefix)
      PREFIX="$2"
      shift 
      shift 
      ;;
    -url|--url)
      URL="$2"
      shift 
      shift 
      ;;
    -dir|--directcory)
      DIRECTORY="$2"
      shift 
      shift 
      ;;
    *)    
      POSITIONAL+=("$1") 
      shift
      ;;
  esac
done

set -- "${POSITIONAL[@]}" # restore positional parameters

echo "FILE EXTENSION  = ${EXTENSION}"
echo "FILE PREFIX     = ${PREFIX}"
echo "SOURCE URL      = ${URL}"
echo "DIRECTORY PATH  = ${DIRECTORY}"

python3 ./fasta-scrapping.py -url ${URL} -prefix ${PREFIX} -ext ${EXTENSION} -dir ${DIRECTORY}