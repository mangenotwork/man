import addContent from './add-content.js';
document.body.innerHTML +='My first Webpack app.<br />';
addContent();
document.body.innerHTML += "<p>JS 动态添加的内容</p>"; // 页面会显示该文本
