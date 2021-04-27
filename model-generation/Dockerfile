FROM python:3.9-slim-buster as base

#
# The intermediary image to install dependencies
#
FROM base as builder

# install pipenv
RUN pip install pipenv

# Make an install dir
RUN mkdir /install
WORKDIR /install

# Copy pipfile and pipenv
COPY Pipfile* /install

# convert pipfile to requirements and install
RUN pipenv lock -r > /requirements.txt \
  && pip install --prefix=/install -r /requirements.txt

#
# The final image we'll deploy
#
FROM base

# copy dependencies into final image
COPY --from=builder /install /usr/local

# HACK: six for some reason won't install via the builder image
RUN pip install six

# Copy necessary files to build image
COPY training /src/training
COPY web /src/web
COPY output/categorize/model/ /data/models/categorize

WORKDIR /src

ENV CAT_MODEL=/data/models/categorize/categorize_tuple.pkl
ENV FLASK_APP=web

EXPOSE 5000
CMD ["gunicorn", "--bind", "0.0.0.0:5000", "web:app"]
