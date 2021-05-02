FROM python:3.7-slim-buster as builder

# install pipenv
RUN pip install pipenv

# Make an install dir
RUN mkdir /install
WORKDIR /install

# Copy pipfile and pipenv
COPY Pipfile* /install

# convert pipfile to requirements and install
RUN pipenv lock -r > requirements.txt \
  && pip install -r requirements.txt

# HACK: six for some reason won't install via the builder image
RUN pip install six

# Copy necessary files to build image
COPY training /src/training
COPY web /src/web
COPY output/categorize/model/ /data/models/categorize

#
# The final image we'll deploy
#
FROM gcr.io/distroless/python3

# copy dependencies into final image
COPY --from=builder /usr/local/lib/python3.7/site-packages/ \
                  /usr/lib/python3.7/.

# copy source data into final image
COPY --from=builder /src /src

# copy data into final image
COPY --from=builder /data /data

WORKDIR /src

ENV CAT_MODEL=/data/models/categorize/categorize_tuple.pkl
ENV FLASK_APP=web

# EXPOSE 5000
CMD ["-m", "gunicorn.app.wsgiapp", "--bind", "0.0.0.0:5000", "web:app"]
