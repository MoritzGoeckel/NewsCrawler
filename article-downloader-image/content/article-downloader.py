import redis
import os
import time as t
from datetime import datetime, date, time, timedelta
import json
from newspaper import Article, Config
from bs4 import BeautifulSoup
import hashlib

def main():
    print('Article downloader version 0.03')

    agtClient, pqClient, lqClient = getRedisConnections()

    while True:
        message = getNextInQueue(lqClient)
        link = json.loads(message[1])
        downloadArticle(url=link['Url'], source=link['Source'], agtClient=agtClient, pqClient=pqClient)

def getRedisConnections():
    agtUrl = os.environ.get('agt-article-redis-url')
    pqUrl = os.environ.get('pq-redis-url')
    lqUrl = os.environ.get('lq-redis-url')

    if agtUrl == '' or pqUrl == '' or lqUrl == '':
        print('Environment variables are not set')

    rp = 6379 #default redis port
    agtClient = redis.Redis(host=agtUrl, port=rp, password='', db=0)
    pqClient = redis.Redis(host=pqUrl, port=rp, password='', db=0)
    lqUrl = redis.Redis(host=lqUrl, port=rp, password='', db=0)
    return agtClient, pqClient, lqUrl

def getNextInQueue(client):
    while True:
        rval = client.blpop(keys='pending', timeout=60)
        if not rval:
            t.sleep(10)
            continue
        return rval

def downloadArticle(url, source, agtClient, pqClient):
    config = Config()
    config.MIN_WORD_COUNT = 100
    article = Article(url=url, config=config)
    try:
        article.download()
        article.parse()
    except:
        return

    if article.meta_lang in ['en']:
        articleHtml = article.html
        soup = BeautifulSoup(articleHtml, 'html.parser')
        description = ''
        try:
            if not description:
                description = soup.find('meta', attrs={'name':'description'}).get('content')
        except:
            if not description:
                description = ''
        try:
            if not description:
                description = soup.find('meta', attrs={'property':'og:description'}).get('content')
        except:
            if not description:
                description = ''
        try:
            if not description:
                description = soup.find('meta', attrs={'property':'twitter:description'}).get('content')
        except:
            if not description:
                description = ''
        a = {
            'Headline': article.title,
            'Description': description,
            'Image': article.top_image,
            'Content': article.text,
            'Source': source,
            'Url': url,
            'Time': str(datetime.now())
        }
        h = hashArticle(a)
        pushed = False
        if not alreadyGotThat(h, agtClient):
            setAlreadyGotThat(h, agtClient)
            data = json.dumps(a)
            pushNewEntry(data, pqClient)
            pushed = True

        pushedMsg = 'agt'
        if pushed:
            pushedMsg = 'new'

        print(pushedMsg, "\t", a)

def hashArticle(a):
    return hashlib.md5(bytes(a['Headline'] + a['Content'] + a['Source'], encoding='utf8')).hexdigest()

def alreadyGotThat(hash, client):
    if client.exists(hash):
        return True
    return False

def setAlreadyGotThat(hash, client):
    expiration = timedelta(hours=72)
    client.set(name=hash, value='seen', ex=expiration)

def pushNewEntry(data, client):
    client.lpush('pending', data)


if __name__ == "__main__":
    main()