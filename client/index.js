var ready = require('./utils/ready');
var body = document.getElementsByTagName('body').item(0);
ready(() => {
  //body.innerText = 'hello world'



  document.getElementById('cpOld').addEventListener('click', function(){
    document.querySelector('.area-old .url').select();
    document.execCommand('copy');
    document.getElementById('cpOld').innerText = 'YEP';
    setTimeout(() => {
      document.getElementById('cpOld').innerText = 'COPY';
    }, 2000);    
  });
  document.getElementById('cpNew').addEventListener('click', function(){
    document.querySelector('.area-new .url').select();
    document.execCommand('copy');
    document.getElementById('cpNew').innerText = 'YEP';
    setTimeout(() => {
      document.getElementById('cpNew').innerText = 'COPY';
    }, 2000);
  });


  document.getElementById('btnGen').addEventListener('click', function(){
    run();
  });
  document.getElementById('randomLink').addEventListener('click', function(){
    chooseRandomExample();
  });
  
  
  chooseRandomExample();
});

function chooseRandomExample(){
  let examples = [
    {usr: 'ninjaprox', repo: 'NVActivityIndicatorView'},
    {usr: 'suzuki-0000', repo: 'HoneycombView'},
    {usr: 'ninjaprox', repo: 'NVActivityIndicatorView'},
    {usr: 'v-braun', repo: 'hero-scrape'},
    {usr: 'Ramotion', repo: 'expanding-collection'},
    {usr: 'v-braun', repo: 'VBPiledView'},
    {usr: 'TBXark', repo: 'TKSwitcherCollection'},
  ];
  var item = examples[Math.floor(Math.random()*examples.length)];
  var usrInput = document.getElementById('usr');
  var repoInput = document.getElementById('repo');
  usrInput.value = item.usr;
  repoInput.value = item.repo;  
  run();
}

function run(){
  var usrInput = document.getElementById('usr');
  var repoInput = document.getElementById('repo');

  let usr = usrInput.value;
  let repo = repoInput.value;
  let url = '/api/' + usr + '/' + repo;
  fetch(url, (err, res) => {
    if(err){
      return;
    }

    document.querySelector('.area-old img').src = res.OgResult.Image;
    document.querySelector('.area-new img').src = res.EndResult.Image;
    
    document.querySelector('.area-old h2').innerText = res.OgResult.Title;
    document.querySelector('.area-new h2').innerText = res.EndResult.Title;

    document.querySelector('.area-old p').innerText = res.OgResult.Description;
    document.querySelector('.area-new p').innerText = res.EndResult.Description;

    var getUrl = window.location;
    document.querySelector('.area-old .url').value = 'https://github.com/' + usr + '/' + repo;
    document.querySelector('.area-new .url').value = getUrl.protocol + '//' + getUrl.hostname + '/' + usr + '/' + repo;

  });    
}


function fetch(url, cb){
  var xhr = new XMLHttpRequest();
  xhr.open('GET', url);
  xhr.onload = function() {
    if (xhr.status === 200) {
      cb(null, JSON.parse(xhr.responseText))
    }
    else{
      cb('error')
    }

  };
  xhr.send();    
}
