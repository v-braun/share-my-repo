var ready = require('./utils/ready');
var body = document.getElementsByTagName('body').item(0);
ready(() => {
  body.innerText = 'hello world'
});
