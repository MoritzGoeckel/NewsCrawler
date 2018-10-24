import os
import string
from collections import Counter
from pymongo import MongoClient
from pymongo.errors import ServerSelectionTimeoutError
from textblob import TextBlob
import nltk
nltk.download('stopwords')
nltk.download('wordnet')
nltk.download('averaged_perceptron_tagger')
from nltk.corpus import stopwords, wordnet
from nltk.stem import WordNetLemmatizer

def main():
    print("Headline Analyzer version 0.01")
    
    english_stopwords = stopwords.words('english')
    lemmatizer = WordNetLemmatizer()
    unicodes2remove = [
        #all kinds of quotes
        u'\u2018', u'\u2019', u'\u201a', u'\u201b', u'\u201c', u'\u201d', u'\u201e', u'\u201f', u'\u2014',
        #all kinds of hyphens
        u'\u002d', u'\u058a', u'\u05be', u'\u1400', u'\u1806', u'\u2010', u'\u2011', u'\u2012', u'\u2013', 
        u'\u2014', u'\u2015', u'\u2e17', u'\u2e1a', u'\u2e3a', u'\u2e3b', u'\u2e40', u'\u301c', u'\u3030',
        u'\u30a0', u'\ufe31', u'\ufe32', u'\ufe58', u'\ufe63', u'\uff0d'
    ]

    mongoClient = getConnection()
    if mongoClient:
        collection_articles = mongoClient.news.articles
        collection_headlines = mongoClient.news.headlines

    headlines = []
    if collection_articles.count_documents({}) > 0:
        articles = collection_articles.find()
        for a in articles:
            h = a['headline']
            if h:
                h = h.translate(str.maketrans('', '', string.punctuation + string.digits))
                for u2r in unicodes2remove:
                    h = h.replace(u2r, '')
                h = h.lower()
                headlines.append(h)
    
    if headlines:
        c = Counter()
        for headline in headlines:
            t = TextBlob(headline)
            ngrams = t.ngrams(3)
            for n in ngrams:
                n = list(n)
                n = ' '.join(n)
                c[n] += 1

        top_ngrams = c.most_common(20) #TODO: just a guess. How many most frequent ngrams do we actually want to store?
        #headlines containing the most frequent ngrams
        headlines_contain_top = set()
        #map of the form {ngram : {headline_1, headline_2, ... , headline_n}} where the headlines (value) contain the ngram (key)
        ngram2headline = {}
        for n_top in top_ngrams:
            headlines_contain_n = set()
            for headline in headlines:
                if n_top[0] in headline:
                    headlines_contain_n.add(headline)
            ngram2headline[n_top[0]] = headlines_contain_n

        #map of the form {ngram : Counter({word1:frequency, word2:frequency})}
        ngram2words = {}
        for ngram, headlines in ngram2headline.items():
            c = Counter()
            for headline in headlines:
                words = list(TextBlob(headline.replace(ngram, '').strip()).words)
                for word in words:
                    word = word.strip()
                    if word not in english_stopwords:
                        pos = nltk.pos_tag([word])[0][1]
                        wn_pos = getWordnetPOS(pos)
                        lemma = ''
                        if wn_pos:
                            lemma = lemmatizer.lemmatize(word, pos=wn_pos)
                        else:
                            lemma = lemmatizer.lemmatize(word)
                        if lemma:
                            c[lemma] += 1
                        else: 
                            print('Lemma was not calculated.')
                            c[word] += 1
            c = c.most_common()
            ngram2words[ngram] = dict(c)

        previous = collection_headlines.find_one_and_replace({'ngram2words': {'$exists':True}}, {'ngram2words': ngram2words}, upsert=True)
        print('updated from', previous)
        print('to', ngram2words)
    
    else:
        print('No headlines found.')

def getWordnetPOS(tag):
    if tag in ['NN', 'NNS', 'NNP', 'NNPS']:
        return wordnet.NOUN
    elif tag in ['VB', 'VBD', 'VBG', 'VBN', 'VBP', 'VBZ']:
        return wordnet.VERB
    elif tag in ['RB', 'RBR', 'RBS']:
        return wordnet.ADV
    elif tag in ['JJ', 'JJR', 'JJS']:
        return wordnet.ADJ
    else:
        return ''

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


if __name__ == "__main__":
    main()