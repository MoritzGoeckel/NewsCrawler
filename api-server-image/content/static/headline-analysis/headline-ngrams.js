var ngrams = [];
var ngramfrequencies = [];

$(function(){
  $.getJSON(url + 'get_headline_ngrams/', function(data){
    data = data['Ngram2Words']
    
    Object.values(data).forEach(h => {
      ngrams.push(h['Ngram']);
      ngramfrequencies.push(h['Count']);
    });

  }).done(function(){

    require(["chart"], function(Chart){

    var ctx = document.getElementById("headline-chart");
    var headlineChart = new Chart(ctx, {
      type: 'horizontalBar',
      data: {
        labels: ngrams.slice(0, 14),
        datasets: [{
          label: '# of Votes',
          data: ngramfrequencies.slice(0, 14),
          backgroundColor: [
            'rgba(77, 77, 0, 0.2)',
            'rgba(102, 102, 0, 0.2)',
            'rgba(128, 128, 0, 0.2)',
            'rgba(153, 153, 0, 0.2)',
            'rgba(179, 179, 0, 0.2)',
            'rgba(204, 204, 0, 0.2)',
            'rgba(230, 230, 0, 0.2)',
            'rgba(255, 255, 0, 0.2)',
            'rgba(255, 255, 0, 0.2)',
            'rgba(255, 255, 0, 0.2)',
            'rgba(255, 255, 0, 0.2)',
            'rgba(255, 255, 0, 0.2)',
            'rgba(255, 255, 0, 0.2)',
            'rgba(255, 255, 0, 0.2)'

          ],
          borderColor: [
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)',
            'rgba(77, 77, 0, 1)'
          ],
          borderWidth: 1
        }]
      },
      options: {
        //responsive: false,
        scales: {
          yAxes: [{
            ticks: {
              fontColor: 'white',
              beginAtZero:true
            }
          }]
        }
      }
    });
    });
    })
  });