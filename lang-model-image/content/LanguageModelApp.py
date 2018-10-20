import os
from pymongo import MongoClient
from pymongo.errors import ServerSelectionTimeoutError
from NgramLanguageModel import Model

def main():
    print('Language Model version 0.01')

    mongoClient = getConnection()
    if mongoClient:
        collection_entropy = mongoClient.news.entropy
        collection_articles = mongoClient.news.articles

        e = collection_entropy.find().sort({'ArticleDateTime':-1}).limit(1)
        latest = e['ArticleDateTime']

        articles = collection_articles.find({'DateTime':{'$gt':latest}})
        model = Model(n=3)
        model.read_frequencies(path='frequencies/reuters_adjusted_freq') #TODO: should this path be part of the env-variables?
        for article in articles:
            article_content = article['Content']
            article_date_time = article['DateTime']
            article_url = article['Url']
            article_id = article['_id']
            pp = model.perplexity(article_content)
            collection_entropy.insert_one({'article_id': article_id,
                                           'article_DateTime': article_date_time,
                                           'article_Url': article_url,
                                           'article_perplexity': pp})


def getConnection():
    mongoUrl = os.environ.get('mongo-url')
    if not mongoUrl:
        print('Environment variables are not set')
    print('mongo url: ', mongoUrl)

    mongoPw = os.environ.get('mongo-pw')
    mongoUser = os.environ.get('mongo-user')
    print('mongo credentials: ', mongoUser, mongoPw)
    print('Connecting to mongo...')

    mp = 27017 #default mongo port
    client = MongoClient(host=mongoUrl, port=mp, username=mongoUser, password=mongoPw)

    try:
        client.server_info()
    except ServerSelectionTimeoutError as err:
        print(err)
        return

    print('Mongo connection established')
    return client