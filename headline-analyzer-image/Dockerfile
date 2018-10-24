FROM python:3

WORKDIR /usr/src/app
COPY ./content .
RUN pip install --no-cache -r requirements.txt
RUN python -m textblob.download_corpora lite


CMD ["python", "./headline-analyzer.py"]

LABEL name="headline-analyzer"