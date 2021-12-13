#!/bin/bash

# grab a gzipped fasta 
POSITIONAL=()
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -fin|--input-file)
      INPUT_FASTA="$2"
      shift # past argument
      shift # past value
      ;;
    -k|--k)
      K="$2"
      shift # past argument
      shift # past value
      ;;
    -mc|--max-cnt)
      MAX_CNT="$2"
      shift # past argument
      shift # past value
      ;;
    -dout|--output-dir)
      OUTPUT_DIR="$2"
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

NAME=$(basename -- "$INPUT_FASTA")
NAME="${NAME%.*}"

mkdir ${OUTPUT_DIR} &> /dev/null
mkdir ${OUTPUT_DIR}/${NAME} &> /dev/null
chmod 755 ${OUTPUT_DIR}/${NAME}  &> /dev/null

# unzip original fasta to std ands get header 
GENOME=$(gzip -d -c ${INPUT_FASTA} | awk 'NR==1{sub(/Chr.*/, ""); print substr($1,2,length($1)-2);}') 
# echo "fasta header is ${GENOME}"

# unzip original fasta to std and split to separate files
cd ${OUTPUT_DIR}/${NAME} 
# echo ${PWD}
gzip -d -c ${INPUT_FASTA} | awk -F '>' 'BEGIN {n_seq=1} /^>/ {F=sprintf("./%d.fasta", n_seq); print > F; n_seq++; next;} {print > F; close(F)}'

# echo "output directory ${OUTPUT_DIR}"
# echo "output directory ${PWD}"

# kmc-ado each file:
OUTPUT_FILE=${OUTPUT_DIR}/${NAME}.json
touch ${OUTPUT_FILE} &> /dev/null

# echo "output directory ${OUTPUT_DIR}"
# echo "output directory ${PWD}"
# echo "output directory ${OUTPUT_FILE}"
printf $'{\n' >> ${OUTPUT_FILE}
printf $'\tgenome: \"' >> ${OUTPUT_FILE} && printf ${GENOME} >> ${OUTPUT_FILE} && printf '\",\n' >> ${OUTPUT_FILE}
printf $'\tsequences: [\n' >> ${OUTPUT_FILE}
for seqfile in ${OUTPUT_DIR}/${NAME}/*.fasta; do
    chmod 755 ${seqfile} 
    name=$(cat ${seqfile} | awk 'NR==1{print substr($1,2);}')
    printf $'\t\t{\n\t\t\tname: \"' >> ${OUTPUT_FILE} && printf ${name} >> ${OUTPUT_FILE} && printf $'\",\n' >> ${OUTPUT_FILE} &&
    printf $'\t\t\tdata: \"' >> ${OUTPUT_FILE} &&
    sudo docker exec kmc_test ./kmc -v -k${K} -cs${MAX_CNT} -fa ${seqfile} ${seqfile}_db . &> /dev/null &&
    sudo docker exec kmc_test ./kmc_dump ${seqfile}_db ${seqfile}_db_dumped &> /dev/null &&
    cat ${seqfile}_db_dumped | awk '{gsub(/\t/,":")}1' | awk '{printf "%s,",$0} END {print ""}' >> ${OUTPUT_FILE} &&
    printf $'\"\n\t\t},\n' >> ${OUTPUT_FILE} 
done
printf $'\t]\n' >> ${OUTPUT_FILE}
printf $'}' >> ${OUTPUT_FILE}

sudo rm -r ${OUTPUT_DIR}/${NAME} &> /dev/null

cat ${OUTPUT_FILE}
sudo rm ${OUTPUT_FILE} &> /dev/null