// Redo another page with the error?

function replaceTag(filt, id, includeFail) {

    var items = data.filter(function(post) {
        return filt(post);
    }).map(function(post) {
        var ad = new Date(post.Date);
        
        var stDate = `${ad.getDate()}/${ad.getMonth()+1}/${ad.getFullYear()} ${ad.getHours()}:${ad.getMinutes()}`;   
        var head = `<div class="demo-card-wide mdl-card mdl-shadow--2dp">
              <div class="mdl-card__title">
                <h1 title="${post.Name}" class="mdl-card__title-text">${post.Item.Title}</h1></br>
              </div>
                <h8>(${stDate})</h8></br>`;
        if (includeFail) {
            head = `${head}<h3> ${post.FailReason}</h3>`;
        }
        var body = `
              <div align="left" class="mdl-card__supporting-text">${post.Item.Description}</div>
              <div class="mdl-card__actions mdl-card--border">
                <a title='${post.Message}' target='_blank' href='${post.Link}' 
                   class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect">Link</a>
              </div>
            </div>
            </br>`;
        return `${head}${body}`
    });
    var allHTML = "";
    items.forEach(function(l) {
        allHTML = `${allHTML}${l}`;
    });
    var elem = document.querySelector(id);
    if (elem != null) {
        elem.innerHTML = allHTML;
    }
}

let state = "compiled and loaded"
console.log(`dynamically loaded ${state}`)
console.log(data);
var notFailed = function(post) {
    return !post.Failed;
};
replaceTag(notFailed, '#list', false);
var failed = function(post) {
    return post.Failed;
};
replaceTag(failed, '#error', true);



export
default {}
