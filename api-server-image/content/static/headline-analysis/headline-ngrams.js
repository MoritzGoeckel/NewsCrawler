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
            'rgba(255, 99, 132, 0.2)',
            'rgba(54, 162, 235, 0.2)',
            'rgba(255, 206, 86, 0.2)',
            'rgba(75, 192, 192, 0.2)',
            'rgba(153, 102, 255, 0.2)',
            'rgba(255, 159, 64, 0.2)',
            'rgba(255, 159, 64, 0.2)',
            'rgba(255, 159, 64, 0.2)',
            'rgba(255, 159, 64, 0.2)',
            'rgba(255, 159, 64, 0.2)',
            'rgba(255, 159, 64, 0.2)',
            'rgba(255, 159, 64, 0.2)',
            'rgba(255, 159, 64, 0.2)',
            'rgba(255, 159, 64, 0.2)'

          ],
          borderColor: [
            'rgba(255,99,132,1)',
            'rgba(54, 162, 235, 1)',
            'rgba(255, 206, 86, 1)',
            'rgba(75, 192, 192, 1)',
            'rgba(153, 102, 255, 1)',
            'rgba(255, 159, 64, 1)',
            'rgba(255, 159, 64, 1)',
            'rgba(255, 159, 64, 1)',
            'rgba(255, 159, 64, 1)',
            'rgba(255, 159, 64, 1)',
            'rgba(255, 159, 64, 1)',
            'rgba(255, 159, 64, 1)',
            'rgba(255, 159, 64, 1)',
            'rgba(255, 159, 64, 1)'
          ],
          borderWidth: 1
        }]
      },
      options: {
        //responsive: false,
        scales: {
          yAxes: [{
            ticks: {
              beginAtZero:true
            }
          }]
        }
      }
    });
    });
    })
  });