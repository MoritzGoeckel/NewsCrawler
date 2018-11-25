var draw_wordcloud = function (clickableWords=false) {

  require(['d3', 'cloud', 'chart'], function (d3, cloud, Chart) {
    const wordcloudContainer = "#wordcloud";
    const text = "text";
    const size = "size";
    //var words = [{ "text": "some", "size": 0 }, { "text": "random", "size": 0 }, { "text": "words", "size": 0 }, { "text": "used", "size": 0 },
    //{ "text": "as", "size": 0 }, { "text": "template", "size": 0 }];
    var words = [];
    var max = 0;
    var min = 0;
    var wordcloudHight = 480;
    var wordcloudWidth = 800;

    var fill = d3.scaleOrdinal(d3.schemeDark2);
    //domain goes from 0 to >100 instead of 10 to 100 in order to prevent the very dark colors at the end of the scales
    var fill2 = d3.scaleSequential(d3.interpolateOranges).domain([0, 110]);//TODO: needs to be adjusted if randomSize() is not used for sizing anymore
    //additional color schemes
    //https://github.com/d3/d3-scale-chromatic/blob/master/README.md#schemeCategory10
    var domainScale = d3.scaleLog().domain([min, max]).range([20, 100]);
    //domain must be set dynamically in relation to the incoming value
    var randomSize = function () {
      return 10 + Math.random() * 90;
    }

    var layout = cloud()
      .size([wordcloudWidth, wordcloudHight])
      .words(words)
      .padding(5)
      .rotate(function () { return ~~(Math.random() * 2) * 0; })
      .font("Impact")
      .fontSize(function (d) { return /*domainScale(d.size);*/ randomSize(); })
      .on("end", draw);

    function draw(words) {
      console.log(words);
      d3.select(wordcloudContainer).append("svg")
        .attr("width", layout.size()[0])
        .attr("height", layout.size()[1])
        .append("g")
        .attr("transform", "translate(" + layout.size()[0] / 2 + "," + layout.size()[1] / 2 + ")")
        .selectAll("text")
        .data(words)
        .enter().append("a")
        .attr("xlink:href", "#")
        .on("click", function (d) { 
          if (clickableWords){
            window.alert(d.text); 
            setWord(d.text, this);
          }
        })
        .append("text")
        .style("font-size", function (d) { return d.size + "px"; })
        .style("font-family", "Impact")
        .style("fill", function (d, i) { return fill2(d.size); })
        .attr("text-anchor", "middle")
        .attr("transform", function (d) {
          return "translate(" + [d.x, d.y] + ")rotate(" + d.rotate + ")";
        })
        .text(function (d) { return d.text; });
    }

    $.getJSON(url + "get_word_cloud", function( data ) {

      //max = data[0].Score;
      //min = data[data.length - 1].Score;

      $.each(data, function(key, val) {
          var token = val.Word
          var freq = val.Score
          words.push({'text':token, 'size':freq})
      });

      layout.start();
    });

    //layout.start();

  });
}

function setWord(index, element)
{
    var word = index;
    //var word = words[index];
    console.log(word);

    $("#articleimages-container").show('slide', { direction: 'left' }, 300);
    
    $(".word").removeClass("activeWord");
    $(element).toggleClass("activeWord");

    //$("#articles").hide();
    //$("#statistics").hide();

    //$("#leftContainer").removeClass("col-md-offset-3");

    //statistic
    $.getJSON(url + "word_statistics/" + word, function( data ) {
        console.log(url + "word_statistics/" + word);
        console.log(data);

        var xData = [];
        var yData = [];

        i = 0;
        $.each( data, function(key, val ) {
            xData.push(val.Date.substr(0, 10));
            yData.push({ y: val.Count, x: Date.parse(val.Date) });
        });

        $("#article-counter-container").show('slide', { direction: 'left' }, 300);

        var ctx = document.getElementById("chartCanvas").getContext("2d");

        if (typeof theLineChart !== 'undefined') {
            theLineChart.destroy();
        }

        theLineChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: xData,
                datasets: [{
                    label: 'Artikel',
                    data: yData,
                    lineTension: 0.2,
                    borderColor: 'rgba(143, 143, 36, 1)',
                    backgroundColor: 'rgba(143, 143, 36, 0.3)'
                }],
                options: {
                    scales: {
                        xAxes: [{
                            type: 'linear',
                            position: 'bottom',
                        }]
                    }
                }
            },
            options: {
                legend: {
                    display: false
                }
            }
        });

        //onResize();
    });

    //Articles
    $.getJSON(url + "elastic_articles/" + word, function (data) {
        console.log(url + "elastic_articles/" + word);
        console.log(data);

        var bilderHTML = "";

        count = 0;
        var output = "";
        $.each(data, function (key, val) {
            if (count < 400)
                output += "<a href='" + val.url + "'>" + shortStr(val.headline, 70) + "</a> <span class='sourceName'>[" + shortStr(val.source, 20) + "]</span><br />";
                
            if (count < 9)
                bilderHTML += "<a target='_blank' href='"+val.url+"'><img src='" + val.image + "' /></a>";
            
            count++;
        });

        $("#related-articles-container").html(output);
        $("#related-articles-container").show('slide', { direction: 'right' }, 300);
        $("#articleimages-container").html(bilderHTML);

        //onResize();
    });
}

function shortStr(string, number)
{
    if(string.length > number)
        string = string.substr(0, number - 4) + " ...";

    return string;
}