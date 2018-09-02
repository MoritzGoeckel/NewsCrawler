//TODO: CROSS ORIGIN RESOURCE SHARING IS NOT WORKING YET

var url = 'http://192.168.99.100:30003/';

var maxWordSize = 50;
var minWordSize = 10;
var uppercaseThreshold = 0.7;

var words;

function onResize()
{
    //$("#bilder").height($("#articles").position().top + $("#articles").height() - $("#bilder").position().top + 10);
}

$(window).resize(function () {
    onResize();
});

$(function(){
    $(".artikelSection").hide();
    $(".verlaufSection").hide();
    $(".bilderSection").hide();

    $.getJSON(url + "get_word_cloud", function( data ) {
        console.log(url + "get_word_cloud");
        console.log(data);

        var max = data[0].Score;
        var min = data[data.length - 1].Score;

        words = data;

        var output = "";
        var index = 0;
        $.each( data, function(key, val ) {
            var proportionalSize = ((val.Score - min) / (max - min));
            var square = (proportionalSize + 1) * (proportionalSize + 1);

            output += "<a href='#' onclick='setWord(" + index + ", this);' class='word'><span style='font-size: " + ((proportionalSize * (maxWordSize - minWordSize)) + minWordSize) + "px; "
                + "line-height: " + (square < 1.3 ? square : 1.3) + ";"
                + (proportionalSize > uppercaseThreshold ? "text-transform: uppercase;" : "") + "' >"
                + val.Word + "</span></a> ";

            index++;
        });

        $("#words").html(output);
    });

    setInterval(function () { onResize(); }, 300);
});

function setWord(index, element)
{
    var word = words[index];
    console.log(word);

    $(".bilderSection").show('slide', { direction: 'left' }, 300);
    
    $(".word").removeClass("activeWord");
    $(element).toggleClass("activeWord");

    //$("#articles").hide();
    //$("#statistics").hide();

    //$("#leftContainer").removeClass("col-md-offset-3");

    //statistic
    $.getJSON(url + "word_statistics/" + word.Word, function( data ) {
        console.log(url + "word_statistics/" + word.Word);
        console.log(data);

        var xData = [];
        var yData = [];

        i = 0;
        $.each( data, function(key, val ) {
            xData.push(val.Date.substr(0, 10));
            yData.push({ y: val.Count, x: Date.parse(val.Date) });
        });

        $(".verlaufSection").show('slide', { direction: 'left' }, 300);

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
                    borderColor: '#0000d2',
                    backgroundColor: 'rgba(0,0,0,0.1)'
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

        onResize();
    });

    //Articles
    $.getJSON(url + "elastic_articles/" + word.Word, function (data) {
        console.log(url + "elastic_articles/" + word.Word);
        console.log(data);

        var bilderHTML = "";

        count = 0;
        var output = "";
        $.each(data, function (key, val) {
            if (count < 400)
                output += "<a href='" + val.url + "'>" + shortStr(val.headline, 70) + "</a> <span class='sourceName'>[" + shortStr(val.source, 20) + "]</span><br />";
                
            if (count < 7)
                bilderHTML += "<a target='_blank' href='"+val.url+"'><img src='" + val.image + "' /></a>";
            
            count++;
        });

        $("#articles").html(output);
        $(".artikelSection").show('slide', { direction: 'right' }, 300);
        $("#bilder").html(bilderHTML);

        onResize();
    });
}

function shortStr(string, number)
{
    if(string.length > number)
        string = string.substr(0, number - 4) + " ...";

    return string;
}