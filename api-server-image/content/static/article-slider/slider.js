      //get high entropy article
      var article_headline = '';
      var article_description = '';
      var imagesrc = '';
      var articlelink = '';
      var carouselItemContainerClass = 'carousel-inner';
      var carouselItem;
      var d = new Date();
      d.setMinutes(d.getMinutes()-30);

      $(function(){
        $.getJSON(url + `get_high_entropy_article_since/${d.toISOString().split('.')[0]+'Z'}`, function(articles){
        console.log("High Entropy-Article: ", articles)
        
        first = articles[0]
        article_headline = first['Headline'];
        article_description = first['Description'];
        imagesrc = first['Image'];
        articlelink = first['Url'];
        carouselItem = `
        <div class="carousel-item active">
          <a target="_blank" href="${articlelink}">
            <img class="d-block w-100 image-overlay" src="${imagesrc}" alt="slide">
          </a>
          <div class="carousel-caption d-none d-md-block">
            <h5>${article_headline}</h5>
              <p>${article_description}</p> 
            </div>
          </div>
        `;
        $(`.${carouselItemContainerClass}`).append(carouselItem);

        Object.values(articles.slice(1,10)).forEach(a => {
          article_headline = a['Headline'];
          article_description = a['Description'];
          imagesrc = a['Image'];
          articlelink = a['Url'];

          carouselItem = `
          <div class="carousel-item">
            <a target="_blank" href="${articlelink}">
              <img class="d-block w-100 image-overlay" src="${imagesrc}" alt="slide">
            </a>
            <div class="carousel-caption d-none d-md-block">
              <h5>${article_headline}</h5>
              <p>${article_description}</p> 
            </div>
          </div>
          `;
       
          $(`.${carouselItemContainerClass}`).append(carouselItem);
        });
        })
      });