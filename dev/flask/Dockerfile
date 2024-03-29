# builder base image
FROM python:3.7.6-slim-buster AS compile-image

# install dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends gcc

# virtualenv
ENV VIRTUAL_ENV=/opt/venv
RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

# add and install requirements
RUN pip install --upgrade pip
COPY ./requirements.txt .
RUN pip install -r requirements.txt

# build-image
FROM python:3.7.6-slim-buster AS runtime-image

# install nc
RUN apt-get update && \
    apt-get install -y --no-install-recommends netcat-openbsd

# copy Python dependencies from build image
COPY --from=compile-image /opt/venv /opt/venv

# set working directory
WORKDIR /usr/src/app

# add user
RUN addgroup --system user && adduser --system --no-create-home --group user
RUN chown -R user:user /usr/src/app && chmod -R 755 /usr/src/app
USER 101

# add app
COPY . /usr/src/app

# set environment variables & copy code from compile image
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1
ENV FLASK_ENV production
ENV APP_SETTINGS project.config.ProductionConfig
ENV PATH="/opt/venv/bin:$PATH"
ARG SECRET_KEY=my_precious
ENV SECRET_KEY $SECRET_KEY

# run server
CMD gunicorn --bind 0.0.0.0:5000 manage:app
