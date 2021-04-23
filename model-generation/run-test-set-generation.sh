#!/bin/bash

# TODO: make this not require TF + Jupyter, maybe just Jupyter
IMAGE_NAME=groceryspend.io/parse-pdf

# check that our custom image is built
# TODO: make this an image we manage and download from ECR or Docker Hub
docker build -t $IMAGE_NAME .

# create output dirs
mkdir -p $PWD/output/categorize/training

PYTHON_ARGS=(
  "/src/training/categorize/extract_pdf_text.py"
  "/data/training/categorize/hyad/2018-Store-Brand-Catalog.pdf"
  "/output/categorize/training/hyad.csv"
)

# run the local image with the following mounts
# - data to get access to the data
# - output to output the file
docker run -it \
  -v "$(pwd)/data:/data:ro" \
  -v "$PWD/training:/src/training:ro" \
  -v "$PWD/output:/output:rw" \
  -e "PYTHONPATH=/src" \
  $IMAGE_NAME \
  python "${PYTHON_ARGS[@]}"
  
  

