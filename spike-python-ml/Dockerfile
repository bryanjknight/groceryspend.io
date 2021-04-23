FROM tensorflow/tensorflow:latest-jupyter

# install magick wand
RUN apt-get update && \
  apt-get install -y \
    libmagickwand-dev

# Allow read on PDFs
RUN sed -i "s/rights=\"none\" pattern=\"PDF\"/rights=\"read\" pattern=\"PDF\"/g" /etc/ImageMagick-6/policy.xml

RUN pip install \
  pdfplumber \
  pandas \
  regex