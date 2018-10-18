from flask import Flask, request, jsonify
from NgramLanguageModel import Model

app = Flask(__name__)

@app.route('/entropy', methods=['POST'])
def calculate_information():
    if request.is_json:
        articles = request.json
        frequencies_path_key = 'frequencies_path'
        frequencies_path_val = articles[frequencies_path_key]
        model = Model(n=3)
        if frequencies_path_val is not None:
            model.read_frequencies(frequencies_path_val)
        else:
            model.read_frequencies()
        information_comparison = {}
        for article_id, article in articles['articles'].items():
            information_comparison[article_id] = model.perplexity(article)
        return jsonify(information_comparison)


if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0')