FROM python:3

WORKDIR /usr/src/app
COPY ./content/ .
RUN pip install --no-cache-dir -r requirements.txt

CMD ["python", "./article-downloader.py"]

LABEL name="downloader"
