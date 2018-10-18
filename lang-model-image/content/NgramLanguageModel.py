from nltk.tokenize import TweetTokenizer
import numpy as np
from collections import Counter
import ujson


class Model(object):

    __slots__= 'n', 'ngram_frequencies', 'tokenizer', 'padding_start', 'padding_end', \
               'ngram_counts', 'nc_values', 'n_tokens', 'unigram_count'

    def __init__(self, n =2, padding_start='<$>', padding_end='</$>'):
        self.n = n
        self.ngram_frequencies = {}
        self.tokenizer = TweetTokenizer()
        self.padding_start = padding_start
        self.padding_end = padding_end
        self.ngram_counts = None #the ngram frequencies of the specific ngram used
        self.nc_values = None #list of all possible frequencies
        self.n_tokens = None #the sum of all values
        self.unigram_count = None #number of unigrams (from the original ngram counts)


    def read_frequencies(self, path='reuters_freq'):
        """
        :param path: the path must be given without the suffix _n1, _n2, _n3
        indicating the actual number n for the ngram frequencies
        :return:
        """
        path = path[:-5] if path.endswith('.json') else path
        with open(path+'_n'+str(self.n)+'.json', mode='r') as f:
            #reader = csv.reader(f, delimiter='\t')
            self.ngram_frequencies['n'+str(self.n)] = Counter(ujson.load(f))
        self.ngram_counts = self.ngram_frequencies['n' + str(self.n)]
        self.nc_values = list(self.ngram_counts.values())
        self.n_tokens = sum(self.nc_values)
        self.unigram_count = self.nc_values.count(2) #2 is used because of the adjusted counts

    def _ngrams(self, tokens):
        """
        helper method
        get ngrams from list of tokens
        :param tokens:
        :return:
        """
        wl = [self.padding_start] * (self.n - 1)  # n-1 starting symbols
        wl += tokens
        wl.append(self.padding_end)
        return zip(*[wl[i:] for i in range(self.n)])

    def mle(self, sequence):
        p = 0
        for ngram in self._ngrams(self.tokenizer.tokenize(sequence)):
            ngram_frequency = self.ngram_counts[ngram]
            # ngram probability from good turing equivalence classes
            if ngram_frequency > 0:
                p += np.log2(self.nc_values.count(ngram_frequency) * ngram_frequency)\
                     - np.log2(self.n_tokens)
            else:
                p += np.log2(self.unigram_count) - np.log2(self.n_tokens)
        print('mle: ', p)
        return p

    def cross_entropy(self, sequence):
        return -(1/len(list(self.tokenizer.tokenize(sequence)))) * self.mle(sequence)

    def perplexity(self, sequence):
        ce = self.cross_entropy(sequence)
        print('cross entropy: ', ce)
        return pow(2, ce)