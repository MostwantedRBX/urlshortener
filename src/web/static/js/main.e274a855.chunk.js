(this["webpackJsonpweb-src"]=this["webpackJsonpweb-src"]||[]).push([[0],{11:function(e,t,n){},14:function(e,t,n){"use strict";n.r(t);var c=n(1),r=n.n(c),s=n(5),o=n.n(s),a=(n(11),n(4)),i=n.n(a),l=n(6),u=n(2),h=n(0),j=function(){return Object(h.jsx)("header",{className:"header",children:Object(h.jsx)("h2",{children:"URL Shortener"})})},b=function(e){var t=e.onShorten,n=Object(c.useState)(""),r=Object(u.a)(n,2),s=r[0],o=r[1];return Object(h.jsxs)("form",{className:"input_form",onSubmit:function(e){e.preventDefault(),!s||s.length<5?alert("Please enter a URL!"):(console.log("On Submit"),t({url:s}),o(""))},children:[Object(h.jsx)("div",{className:"input_div",children:Object(h.jsx)("input",{className:"input_field",type:"text",placeholder:"https://www.google.com/",value:s,onChange:function(e){return o(e.target.value)}})}),Object(h.jsx)("input",{className:"btn btn-block",type:"submit",value:"Shorten URL"})]})};var p=function(){var e=Object(c.useState)(""),t=Object(u.a)(e,2),n=t[0],r=t[1];function s(){return(s=Object(l.a)(i.a.mark((function e(t){var n,c;return i.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return console.log(t.url),console.log("On shortenUrl"),e.next=4,fetch("http://localhost:8080/links/put/",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify(t)});case 4:return n=e.sent,e.next=7,n.json();case 7:c=e.sent,r(c.url);case 9:case"end":return e.stop()}}),e)})))).apply(this,arguments)}return Object(h.jsxs)("div",{className:"container",children:[Object(h.jsx)(j,{}),console.log("Rendered"),Object(h.jsx)(b,{onShorten:function(e){return s.apply(this,arguments)}}),Object(h.jsx)("p",{className:"urlP",children:n||""})]})},d=function(e){e&&e instanceof Function&&n.e(3).then(n.bind(null,15)).then((function(t){var n=t.getCLS,c=t.getFID,r=t.getFCP,s=t.getLCP,o=t.getTTFB;n(e),c(e),r(e),s(e),o(e)}))};o.a.render(Object(h.jsx)(r.a.StrictMode,{children:Object(h.jsx)(p,{})}),document.getElementById("root")),d()}},[[14,1,2]]]);
//# sourceMappingURL=main.e274a855.chunk.js.map