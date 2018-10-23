FROM python:3

WORKDIR /usr/src/app
COPY ./content .
RUN pip install --no-cache -r requirements.txt

CMD ["python", "./LanguageModelApp.py"]

LABEL name="language-model"