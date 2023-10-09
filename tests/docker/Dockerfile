FROM python:3.12-slim

RUN mkdir /workspace

WORKDIR /workspace

COPY tests tests

WORKDIR /workspace/tests

RUN pip install -r requirements.txt

ENTRYPOINT [ "python3", "-m", "pytest"]
