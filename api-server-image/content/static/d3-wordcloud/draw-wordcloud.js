require(['d3', 'cloud'], function(d3, cloud){
  const wordcloudContainer = "#wordcloud";
  const text = "text";
  const size = "size";
  var words = [];
  var max = 0;
  var min = 0;
  var wordcloudHight = 500;
  var wordcloudWidth = 800;

  var fill = d3.scaleOrdinal(d3.schemeDark2);
  //additional color schemes
  //https://github.com/d3/d3-scale-chromatic/blob/master/README.md#schemeCategory10
  var domainScale = d3.scaleLog().domain([min, max]).range([20,100]);
  //domain must be set dynamically in relation to the incoming value
  var randomSize = function(){
    return 10 + Math.random() * 90;
  }

  var layout = cloud()
    .size([wordcloudWidth, wordcloudHight])
    .words(words)
    .padding(5)
    .rotate(function() { return ~~(Math.random() * 2) * 0; })
    .font("Impact")
    .fontSize(function(d) { return /*domainScale(d.size);*/ randomSize(); })
    .on("end", draw);

  function draw(words) {
      console.log("draw..")
      console.log(words)
      d3.select(wordcloudContainer).append("svg")
          .attr("width", layout.size()[0])
          .attr("height", layout.size()[1])
        .append("g")
          .attr("transform", "translate(" + layout.size()[0] / 2 + "," + layout.size()[1] / 2 + ")")
        .selectAll("text")
          .data(words)
        .enter().append("text")
          .style("font-size", function(d) { return d.size + "px"; })
          .style("font-family", "Impact")
          .style("fill", function(d,i){return fill(i); })
          .attr("text-anchor", "middle")
          .attr("transform", function(d) {
            return "translate(" + [d.x, d.y] + ")rotate(" + d.rotate + ")";
          })
          .text(function(d) { return d.text; });
  }

  $.getJSON(url + "get_word_cloud", function( data ) {
    console.log(url + "get_word_cloud");
  
    max = data[0].Score;
    min = data[data.length - 1].Score;
  
    $.each(data, function(key, val) {
        var token = val.Word
        var freq = val.Score
        console.log(token, freq)
        words.push({'text':token, 'size':freq})
    });
   
    console.log("max: ", max);
    console.log("min: ", min);
    console.log("words: ");
    console.log(words);
    layout.start();
  });
  

});