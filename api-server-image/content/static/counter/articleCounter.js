var thisYear = '';
var thisMonth = '';
var thisDay = '';
var dates = [];
var yearCount;
var monthCount;
var dayCount;
var counterIds = [
  'article-counter-day',
  'article-counter-month', 
  'article-counter-year'
];
var AJAX = [];
var options = {
        useEasing: true, 
        useGrouping: false, 
        separator: ',', 
        decimal: '.', 
};
var d = new Date();

var getCount = function(from){
  return $.getJSON(url + `get_article_count_since/${from}`);
}

var count = function(htmlId, c, o=options){
  var counter = new CountUp(htmlId, 0, c, 0, 2.5, o);
  if (!counter.error){
    counter.start();
  } else {
    console.error(counter.error);
  }
}


d.setHours(1, 0, 0, 0); //for some reason 1 needs to be set to get 00:00:00 when converting to rfc3339
thisDay = d.toISOString().split('.')[0]+'Z';

d.setDate(1);
thisMonth = d.toISOString().split('.')[0]+'Z';

d.setMonth(0, 1);
thisYear = d.toISOString().split('.')[0]+'Z';

dates.push(thisDay);
dates.push(thisMonth);
dates.push(thisYear);
for(var i=0; i < dates.length; i++){
  AJAX.push(getCount(dates[i]));
}
$.when.apply($, AJAX).done(function(){
  var counts = [];
  for(var i=0; i < arguments.length; i++){
    counts.push(arguments[i][0]);
  }
  for(var i=0; i < counts.length; i++){
    count(counterIds[i], counts[i]);
  }

  require(["chart"], function(Chart){
  // For a pie chart
  var ctx2 = document.getElementById("counter-pie");
  var counterPieChart = new Chart(ctx2,{
    type: 'pie',
    data: data = {
      datasets: [{
        data: counts,
        backgroundColor: [
          'rgba(255, 99, 132, 0.2)',
          'rgba(54, 162, 235, 0.2)',
          'rgba(255, 206, 86, 0.2)'
        ],
        borderColor: [
          'rgba(75, 192, 192, 1)',
          'rgba(153, 102, 255, 1)',
          'rgba(255, 159, 64, 1)'
        ]
      }],
    // These labels appear in the legend and in the tooltips when hovering different arcs
    labels: [
      'Day count',
      'Month count',
      'Year count'
    ]
  }//,
  //options: options
});
});
});