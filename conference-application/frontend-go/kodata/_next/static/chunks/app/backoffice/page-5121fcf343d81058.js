(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[280],{629:function(e,t,n){Promise.resolve().then(n.bind(n,8950))},8950:function(e,t,n){"use strict";let s;n.r(t),n.d(t,{default:function(){return q}});var a=n(7437),c=n(9617),i=n.n(c),o=n(2265),r=n(52),l=n(7157),d=n.n(l),p=function(e){let{title:t,author:n,id:s,status:c,approved:i,email:o,description:l,actionHandler:p}=e,h=(e,t,n)=>{p(t,e,n)};return(0,a.jsxs)("div",{className:"".concat(d().ProposalItem,"  ").concat("PENDING"===c?d().pending:"","   ").concat("DECIDED"===c?d().decided:"","  ").concat("REJECT"===c?d().rejected:"","   ").concat("APPROVE"===c?d().approved:"","  ").concat("ARCHIVED"===c?d().archived:""),children:[(0,a.jsxs)("div",{className:"ProposalItem__header",children:[(0,a.jsx)("h4",{children:t}),(0,a.jsxs)("div",{children:[n," ",o]}),(0,a.jsx)("div",{className:d().status,children:c})]}),(0,a.jsx)("div",{className:d().description,children:(0,a.jsx)("p",{className:"p --s",children:l})}),!1,c&&"PENDING"===c&&(0,a.jsxs)("div",{className:d().actions,children:[(0,a.jsx)("div",{children:(0,a.jsx)(r.default,{clickHandler:()=>h(s,c,"APPROVE"),children:"Approve"})}),(0,a.jsx)("div",{children:(0,a.jsx)(r.default,{clickHandler:()=>h(s,c,"REJECT"),children:"Reject"})}),(0,a.jsx)("div",{children:(0,a.jsx)(r.default,{clickHandler:()=>h(s,c,"ARCHIVE"),children:"Archive"})})]}),c&&"DECIDED"===c&&(0,a.jsxs)("div",{className:d().actions,children:[!0===i&&(0,a.jsx)("div",{className:"".concat(d().statusTag,"  ").concat(d().approved),children:"Approved"}),!1===i&&(0,a.jsx)("div",{className:"".concat(d().statusTag,"  ").concat(d().rejected),children:"Rejected"}),(0,a.jsx)("div",{children:(0,a.jsx)(r.default,{clickHandler:()=>h(s,c,"ARCHIVE"),children:"Archive"})})]}),c&&"ARCHIVED"===c&&(0,a.jsx)("div",{className:d().actions,children:(0,a.jsx)("div",{className:"ProposalItem__badge --approved",children:"Archived"})})]})},h=function(){let[e,t]=(0,o.useState)(!1),[n,s]=(0,o.useState)(!1),[c,i]=(0,o.useState)("PENDING"),[r,l]=(0,o.useState)([]),h=(e,n)=>{t(!0),s(!1),console.log("Decision Made ..."),fetch("/api/c4p/proposals/"+e+"/decide/",{method:"POST",body:JSON.stringify({approved:n}),headers:{accept:"application/json"}}).then(e=>e.json()).then(e=>{var n="?status="+c;""==c&&(n=""),m(n),t(!1)}).catch(e=>{t(!1),s(!0)})},u=e=>{t(!0),s(!1),console.log("Archiving Proposal ..."+e),fetch("/api/c4p/proposals/"+e,{method:"DELETE",headers:{accept:"application/json"}}).then(e=>e.json()).then(()=>{var e="?status="+c;""==c&&(e=""),m(e),t(!1)}).catch(e=>{console.log(e),t(!1),s(!0)})};function f(e,t,n){console.log("status: "+e+" - id: "+t+" - action: "+n),"PENDING"==e&&("APPROVE"==n?h(t,!0):h(t,!1)),"ARCHIVE"==n&&u(t)}let m=e=>{console.log("Fetching Proposals ... ("+e+")."),fetch("/api/c4p/proposals/"+e).then(e=>e.json()).then(e=>{l(e),t(!1)}).catch(e=>{console.log(e)})};return(0,o.useEffect)(()=>{t(!0);var e="?status="+c;""==c&&(e=""),m(e)},[c]),(0,a.jsxs)("div",{className:d().ProposalList,children:[(0,a.jsx)("div",{className:d().ProposalList_Filters,children:(0,a.jsxs)("div",{className:d().container,children:[(0,a.jsx)("div",{className:d().filterLabel,children:"Filter By: "}),(0,a.jsx)("div",{className:"".concat(""==c?d().inactive:d().active,"  ").concat(d().filter),onClick:()=>void i(""),children:"All"}),(0,a.jsx)("div",{className:"".concat("PENDING"==c?d().inactive:d().active,"   ").concat(d().filter),onClick:()=>void i("PENDING"),children:"Pending"}),(0,a.jsx)("div",{className:"".concat("DECIDED"==c?d().inactive:d().active,"   ").concat(d().filter),onClick:()=>void i("DECIDED"),children:"Decided"}),(0,a.jsx)("div",{className:"".concat("ARCHIVED"==c?d().inactive:d().active,"   ").concat(d().filter),onClick:()=>void i("ARCHIVED"),children:"Archived"})]})}),(0,a.jsxs)("div",{className:d().ProposalList_Items,children:[r&&r.map((e,t)=>(0,a.jsx)(p,{id:e.id,title:e.title,author:e.author,description:e.description,email:e.email,approved:e.approved,status:e.status.status,actionHandler:f},e.id)),r&&0===r.length&&"PENDING"===c&&(0,a.jsx)("span",{children:"There are no pending proposals."}),r&&0===r.length&&"DECIDED"===c&&(0,a.jsx)("span",{children:"There are no decided proposals."}),r&&0===r.length&&!1===c&&(0,a.jsx)("span",{children:"There are no proposals."})]})]})},u=n(2880),f=n.n(u),m=function(e){let{id:t,title:n,emailTo:s,emailSubject:c,emailBody:i,approved:r}=e,[l,d]=(0,o.useState)(!1),p=()=>{l?d(!1):d(!0)};return(0,a.jsxs)("div",{onClick:()=>p(),className:"".concat(f().NotificationItem,"  ").concat(l?f().open:" "," "),children:[(0,a.jsxs)("div",{className:f().openTag,children:[!l&&(0,a.jsx)(a.Fragment,{children:"Click for details"}),l&&(0,a.jsx)(a.Fragment,{children:"Close"})]}),(0,a.jsxs)("div",{className:f().header,children:[(0,a.jsxs)("h5",{children:[" ",(0,a.jsx)("span",{children:"Proposal:"})," ",n]}),(0,a.jsxs)("div",{className:f().headerStatus,children:[r&&(0,a.jsx)("div",{className:"".concat(f().headerStatusTag,"  ").concat(f().approved," "),children:"Approved"}),!r&&(0,a.jsx)("div",{className:"".concat(f().headerStatusTag,"  ").concat(f().rejected," "),children:"Rejected"})]})]}),(0,a.jsxs)("div",{className:f().description,children:[(0,a.jsxs)("div",{className:f().descriptionTo,children:[(0,a.jsx)("span",{children:"To:"})," ",s]}),(0,a.jsxs)("div",{className:f().descriptionSubject,children:[(0,a.jsx)("span",{children:" Subject:"})," ",c]}),(0,a.jsx)("div",{className:f().descriptionBody,children:(0,a.jsx)("p",{children:i})})]})]})},_=function(){let[e,t]=(0,o.useState)(!1),[n,s]=(0,o.useState)(!1),[c,i]=(0,o.useState)([]),[r,l]=(0,o.useState)(0),d=()=>{fetch("/api/notifications/notifications/").then(e=>e.json()).then(e=>{console.log("Fetching Notifications ..."),i(e),t(!1)}).catch(e=>{console.log(e)})};return(0,o.useEffect)(()=>{let e=setInterval(()=>{t(!0),d()},3e3);return()=>clearInterval(e)},[r]),(0,o.useEffect)(()=>{t(!0),d()},[]),(0,a.jsx)("div",{children:(0,a.jsxs)("div",{className:f().NotificationList,children:[c&&c.map((e,t)=>(0,a.jsx)(m,{id:e.id,title:e.title,emailTo:e.emailTo,emailBody:e.emailBody,emailSubject:e.emailSubject,approved:e.accepted},e.id)),c&&0===c.length&&(0,a.jsx)("span",{children:"There are no notifications."})]})})},v=n(6782),j=n.n(v),N=function(e){let{id:t,type:n,payload:s}=e,[c,i]=(0,o.useState)(!1),r=()=>{c?i(!1):i(!0)};return(0,a.jsxs)("div",{onClick:()=>r(),className:"".concat(j().EventItem,"  ").concat(c?j().open:" "," "),children:[(0,a.jsxs)("div",{className:j().openTag,children:[!c&&(0,a.jsx)(a.Fragment,{children:"Click for details"}),c&&(0,a.jsx)(a.Fragment,{children:"Close"})]}),(0,a.jsx)("div",{className:j().header,children:(0,a.jsxs)("h5",{children:[(0,a.jsxs)("span",{children:["#",t]}),"  ",n]})}),(0,a.jsx)("div",{className:j().description,children:(0,a.jsx)("div",{className:j().codeContainer,children:s})})]})},x=function(){let[e,t]=(0,o.useState)(!1),[n,s]=(0,o.useState)(!1),[c,i]=(0,o.useState)([]),[r,l]=(0,o.useState)(0),d=()=>{fetch("/api/events/").then(e=>e.json()).then(e=>{console.log("Fetching Events ..."),i(e),t(!1)}).catch(e=>{console.log(e)})};return(0,o.useEffect)(()=>{t(!0),d()},[]),(0,o.useEffect)(()=>{let e=setInterval(()=>{t(!0),d()},3e3);return()=>clearInterval(e)},[r]),(0,a.jsx)("div",{children:(0,a.jsxs)("div",{className:j().EventsList,children:[c&&c.map((e,t)=>(0,a.jsx)(N,{id:e.Id,type:e.Type,payload:e.Payload},e.Id)),c&&0===c.length&&(0,a.jsx)("span",{children:"There are no events."})]})})};function g(e){return t=>!!t.type&&t.type.tabsRole===e}let b=g("Tab"),A=g("TabList"),I=g("TabPanel");function E(e,t){return o.Children.map(e,e=>null===e?null:b(e)||A(e)||I(e)?t(e):e.props&&e.props.children&&"object"==typeof e.props.children?(0,o.cloneElement)(e,{...e.props,children:E(e.props.children,t)}):e)}var y=function(){for(var e,t,n=0,s="";n<arguments.length;)(e=arguments[n++])&&(t=function e(t){var n,s,a="";if("string"==typeof t||"number"==typeof t)a+=t;else if("object"==typeof t){if(Array.isArray(t))for(n=0;n<t.length;n++)t[n]&&(s=e(t[n]))&&(a&&(a+=" "),a+=s);else for(n in t)t[n]&&(a&&(a+=" "),a+=n)}return a}(e))&&(s&&(s+=" "),s+=t);return s};function C(e){let t=0;return!function e(t,n){return o.Children.forEach(t,t=>{null!==t&&(b(t)||I(t)?n(t):t.props&&t.props.children&&"object"==typeof t.props.children&&(A(t)&&n(t),e(t.props.children,n)))})}(e,e=>{b(e)&&t++}),t}function T(e){return e&&"getAttribute"in e}function S(e){return T(e)&&e.getAttribute("data-rttab")}function k(e){return T(e)&&"true"===e.getAttribute("aria-disabled")}let D={className:"react-tabs",focus:!1},P=e=>{let t=(0,o.useRef)([]),n=(0,o.useRef)([]),a=(0,o.useRef)();function c(t,n){if(t<0||t>=l())return;let{onSelect:s,selectedIndex:a}=e;s(t,a,n)}function i(e){let t=l();for(let n=e+1;n<t;n++)if(!k(d(n)))return n;for(let t=0;t<e;t++)if(!k(d(t)))return t;return e}function r(e){let t=e;for(;t--;)if(!k(d(t)))return t;for(t=l();t-- >e;)if(!k(d(t)))return t;return e}function l(){let{children:t}=e;return C(t)}function d(e){return t.current[`tabs-${e}`]}function p(e){let t=e.target;do if(h(t)){if(k(t))return;let n=[].slice.call(t.parentNode.children).filter(S).indexOf(t);c(n,e);return}while(null!=(t=t.parentNode))}function h(e){if(!S(e))return!1;let t=e.parentElement;do{if(t===a.current)return!0;if(t.getAttribute("data-rttabs"))break;t=t.parentElement}while(t);return!1}let{children:u,className:f,disabledTabClassName:m,domRef:_,focus:v,forceRenderTabPanel:j,onSelect:N,selectedIndex:x,selectedTabClassName:g,selectedTabPanelClassName:T,environment:P,disableUpDownKeys:R,disableLeftRightKeys:L,...w}={...D,...e};return o.createElement("div",Object.assign({},w,{className:y(f),onClick:p,onKeyDown:function(t){let{direction:n,disableUpDownKeys:s,disableLeftRightKeys:a}=e;if(h(t.target)){let{selectedIndex:o}=e,h=!1,u=!1;("Space"===t.code||32===t.keyCode||"Enter"===t.code||13===t.keyCode)&&(h=!0,u=!1,p(t)),(a||37!==t.keyCode&&"ArrowLeft"!==t.code)&&(s||38!==t.keyCode&&"ArrowUp"!==t.code)?(a||39!==t.keyCode&&"ArrowRight"!==t.code)&&(s||40!==t.keyCode&&"ArrowDown"!==t.code)?35===t.keyCode||"End"===t.code?(o=function(){let e=l();for(;e--;)if(!k(d(e)))return e;return null}(),h=!0,u=!0):(36===t.keyCode||"Home"===t.code)&&(o=function(){let e=l();for(let t=0;t<e;t++)if(!k(d(t)))return t;return null}(),h=!0,u=!0):(o="rtl"===n?r(o):i(o),h=!0,u=!0):(o="rtl"===n?i(o):r(o),h=!0,u=!0),h&&t.preventDefault(),u&&c(o,t)}},ref:e=>{a.current=e,_&&_(e)},"data-rttabs":!0}),function(){let a=0,{children:c,disabledTabClassName:i,focus:r,forceRenderTabPanel:p,selectedIndex:h,selectedTabClassName:u,selectedTabPanelClassName:f,environment:m}=e;n.current=n.current||[];let _=n.current.length-l(),v=(0,o.useId)();for(;_++<0;)n.current.push(`${v}${n.current.length}`);return E(c,e=>{let c=e;if(A(e)){let a=0,l=!1;null==s&&function(e){let t=e||("undefined"!=typeof window?window:void 0);try{s=!!(void 0!==t&&t.document&&t.document.activeElement)}catch(e){s=!1}}(m);let p=m||("undefined"!=typeof window?window:void 0);s&&p&&(l=o.Children.toArray(e.props.children).filter(b).some((e,t)=>p.document.activeElement===d(t))),c=(0,o.cloneElement)(e,{children:E(e.props.children,e=>{let s=`tabs-${a}`,c=h===a,d={tabRef:e=>{t.current[s]=e},id:n.current[a],selected:c,focus:c&&(r||l)};return u&&(d.selectedClassName=u),i&&(d.disabledClassName=i),a++,(0,o.cloneElement)(e,d)})})}else if(I(e)){let t={id:n.current[a],selected:h===a};p&&(t.forceRender=p),f&&(t.selectedClassName=f),a++,c=(0,o.cloneElement)(e,t)}return c})}())};P.propTypes={};let R={defaultFocus:!1,focusTabOnClick:!0,forceRenderTabPanel:!1,selectedIndex:null,defaultIndex:null,environment:null,disableUpDownKeys:!1,disableLeftRightKeys:!1},L=e=>null===e.selectedIndex?1:0,w=(e,t)=>{},F=e=>{let{children:t,defaultFocus:n,defaultIndex:s,focusTabOnClick:a,onSelect:c,...i}={...R,...e},[r,l]=(0,o.useState)(n),[d]=(0,o.useState)(L(i)),[p,h]=(0,o.useState)(1===d?s||0:null);if((0,o.useEffect)(()=>{l(!1)},[]),1===d){let e=C(t);(0,o.useEffect)(()=>{if(null!=p){let t=Math.max(0,e-1);h(Math.min(p,t))}},[e])}w(i,d);let u={...e,...i};return u.focus=r,u.onSelect=(e,t,n)=>{("function"!=typeof c||!1!==c(e,t,n))&&(a&&l(!0),1===d&&h(e))},null!=p&&(u.selectedIndex=p),delete u.defaultFocus,delete u.defaultIndex,delete u.focusTabOnClick,o.createElement(P,u,t)};F.propTypes={},F.tabsRole="Tabs";let H={className:"react-tabs__tab-list"},B=e=>{let{children:t,className:n,...s}={...H,...e};return o.createElement("ul",Object.assign({},s,{className:y(n),role:"tablist"}),t)};B.tabsRole="TabList",B.propTypes={};let O="react-tabs__tab",G={className:O,disabledClassName:`${O}--disabled`,focus:!1,id:null,selected:!1,selectedClassName:`${O}--selected`},K=e=>{let t=(0,o.useRef)(),{children:n,className:s,disabled:a,disabledClassName:c,focus:i,id:r,selected:l,selectedClassName:d,tabIndex:p,tabRef:h,...u}={...G,...e};return(0,o.useEffect)(()=>{l&&i&&t.current.focus()},[l,i]),o.createElement("li",Object.assign({},u,{className:y(s,{[d]:l,[c]:a}),ref:e=>{t.current=e,h&&h(e)},role:"tab",id:`tab${r}`,"aria-selected":l?"true":"false","aria-disabled":a?"true":"false","aria-controls":`panel${r}`,tabIndex:p||(l?"0":null),"data-rttab":!0}),n)};K.propTypes={},K.tabsRole="Tab";let V="react-tabs__tab-panel",$={className:V,forceRender:!1,selectedClassName:`${V}--selected`},Q=e=>{let{children:t,className:n,forceRender:s,id:a,selected:c,selectedClassName:i,...r}={...$,...e};return o.createElement("div",Object.assign({},r,{className:y(n,{[i]:c}),role:"tabpanel",id:`panel${a}`,"aria-labelledby":`tab${a}`}),s||c?t:null)};Q.tabsRole="TabPanel",Q.propTypes={};var U=n(3110),W=n(4902),Y=n.n(W),z=function(e){let{key:t,name:n,version:s,source:c,podName:i,nodeName:o,namespace:r,podIp:l,serviceAccount:d,healthy:p}=e;return(0,a.jsx)("div",{className:Y().ServiceInfo,children:(0,a.jsxs)("div",{children:[(0,a.jsxs)("div",{className:Y().header,children:[(0,a.jsx)("h4",{children:(0,a.jsx)("a",{href:c,target:"_blank",children:n})}),(0,a.jsx)("h5",{children:s})]}),(0,a.jsxs)("div",{className:Y().description,children:[(0,a.jsxs)("div",{className:Y().descriptionItem,children:[(0,a.jsx)("span",{children:"Pod Name:"}),i]}),(0,a.jsxs)("div",{className:Y().descriptionItem,children:[(0,a.jsx)("span",{children:"Node Name:"}),o]}),(0,a.jsxs)("div",{className:Y().descriptionItem,children:[(0,a.jsx)("span",{children:"Pod Namespace:"}),r]}),(0,a.jsxs)("div",{className:Y().descriptionItem,children:[(0,a.jsx)("span",{children:"Pod IP:"}),l]}),(0,a.jsxs)("div",{className:Y().descriptionItem,children:[(0,a.jsx)("span",{children:"Pod Service Account:"}),d]})]}),(0,a.jsxs)("div",{className:"".concat(Y().statusTag,"  ").concat(null!=p&&p?Y().healthy:Y().unhealthy),children:[null!=p&&p&&(0,a.jsx)(a.Fragment,{children:"Healthy"}),null!=p&&!p&&(0,a.jsx)(a.Fragment,{children:"Unhealthy"})]})]})})},M=function(){let[e,t]=(0,o.useState)(!1),[n,s]=(0,o.useState)(""),[c,i]=(0,o.useState)(""),[r,l]=(0,o.useState)(""),[d,p]=(0,o.useState)(""),[h,u]=(0,o.useState)(0),f={name:"FRONTEND",podId:"N/A",podNamespace:"N/A",podNodeName:"N/A",podName:"N/A",podServiceAccount:"N/A",source:"N/A",version:"N/A",healthy:!1},m={name:"AGENDA",podId:"N/A",podNamespace:"N/A",podNodeName:"N/A",podName:"N/A",podServiceAccount:"N/A",source:"N/A",version:"N/A",healthy:!1},_={name:"C4P",podId:"N/A",podNamespace:"N/A",podNodeName:"N/A",podName:"N/A",podServiceAccount:"N/A",source:"N/A",version:"N/A",healthy:!1},v={name:"NOTIFICATIONS",podId:"N/A",podNamespace:"N/A",podNodeName:"N/A",podName:"N/A",serviceAccount:"N/A",source:"N/A",version:"N/A",healthy:!1};(0,o.useEffect)(()=>{let e=setInterval(()=>{t(!0),j(),N(),x(),g()},3e3);return()=>clearInterval(e)},[h]);let j=()=>{t(!0),console.log("Querying service/info"),fetch("/service/info").then(e=>e.json()).then(e=>{e.Healthy=!0,s(e),t(!1)}).catch(e=>{s(f),console.log(e)})},N=()=>{t(!0),console.log("Querying /api/agenda/service/info"),fetch("/api/agenda/service/info").then(e=>e.json()).then(e=>{e.Healthy=!0,l(e),t(!1)}).catch(e=>{l(m),console.log(e)})},x=()=>{t(!0),console.log("Querying /api/c4p/service/info"),fetch("/api/c4p/service/info").then(e=>e.json()).then(e=>{e.Healthy=!0,i(e),t(!1)}).catch(e=>{i(_),console.log(e)})},g=()=>{t(!0),console.log("Querying /api/notifications/service/info"),fetch("/api/notifications/service/info").then(e=>e.json()).then(e=>{e.Healthy=!0,p(e),t(!1)}).catch(e=>{p(v),console.log(e)})};return(0,o.useEffect)(()=>{t(!0),j(),N(),x(),g()},[]),(0,a.jsxs)("div",{className:Y().DebugList,children:[(0,a.jsx)(z,{name:n.name,version:n.version,source:n.source,podIp:n.podIp,podName:n.podName,nodeName:n.podNodeName,namespace:n.podNamespace,serviceAccount:n.podServiceAccount,healthy:n.healthy},"frontend"),(0,a.jsx)(z,{name:c.name,version:c.version,source:c.source,podName:c.podName,podIp:c.podIp,nodeName:c.podNodeName,namespace:c.podNamespace,serviceAccount:c.podServiceAccount,healthy:c.healthy},"c4p"),(0,a.jsx)(z,{name:r.name,version:r.version,source:r.source,podIp:r.podIp,podName:r.podName,nodeName:r.podNodeName,namespace:r.podNamespace,serviceAccount:r.podServiceAccount,healthy:r.healthy},"agenda"),(0,a.jsx)(z,{name:d.name,version:d.version,source:d.source,podIp:d.podIp,podName:d.podName,nodeName:d.podNodeName,namespace:d.podNamespace,serviceAccount:d.podServiceAccount,healthy:d.healthy},"notifications")]})};function q(){let[e,t]=(0,o.useState)(!1),[n,s]=(0,o.useState)(""),c=()=>{t(!0),console.log("Querying /api/features/"),fetch("/api/features/").then(e=>e.json()).then(e=>{console.log("Features Data: "+e),s(e),t(!1)}).catch(e=>{s({}),console.log(e)})};return(0,o.useEffect)(()=>{t(!0),c()},[]),(0,a.jsxs)("main",{className:i().main,children:[(0,a.jsx)("div",{className:"".concat(i().hero," "),children:(0,a.jsx)("div",{className:"grid content noMargin",children:(0,a.jsx)("div",{className:"col full",children:(0,a.jsx)("h3",{children:"Backoffice"})})})}),(0,a.jsx)("div",{className:"".concat(i().BackofficeContent," "),children:(0,a.jsx)("div",{className:"grid content noMargin",children:(0,a.jsx)("div",{className:"col full",children:(0,a.jsx)("div",{className:"".concat(i().tabs," "),children:(0,a.jsxs)(F,{children:[(0,a.jsxs)(B,{children:[(0,a.jsx)(K,{children:"Review Proposals"}),(0,a.jsx)(K,{children:"Agenda Items"}),(0,a.jsx)(K,{children:"Notifications"}),(0,a.jsx)(K,{children:"Events"}),"true"==n.DEBUG_ENABLED&&(0,a.jsx)(K,{children:"Debug"})]}),(0,a.jsx)(Q,{children:(0,a.jsx)(h,{})}),(0,a.jsx)(Q,{children:(0,a.jsx)(U.Z,{admin:"true"})}),(0,a.jsx)(Q,{children:(0,a.jsx)(_,{})}),(0,a.jsx)(Q,{children:(0,a.jsx)(x,{})}),"true"==n.DEBUG_ENABLED&&(0,a.jsx)(Q,{children:(0,a.jsx)(M,{})})]})})})})})]})}},3110:function(e,t,n){"use strict";n.d(t,{Z:function(){return l}});var s=n(7437),a=n(1449),c=n.n(a),i=n(2265),o=n(52),r=function(e){let{key:t,id:n,name:a,day:r,time:l,author:d,description:p,admin:h,handleArchive:u}=e,[f,m]=(0,i.useState)(!1),_=e=>{u(e)},v=()=>{f?m(!1):m(!0)};return(0,s.jsxs)("div",{onClick:()=>v(),className:"".concat(c().agendaItem,"  ").concat(f?c().open:" "," "),children:[(0,s.jsxs)("div",{className:c().openTag,children:[!f&&(0,s.jsx)(s.Fragment,{children:"Click for details"}),f&&(0,s.jsx)(s.Fragment,{children:"Close"})]}),(0,s.jsxs)("div",{className:"AgendaItem__date",children:[(0,s.jsx)("div",{className:"AgendaItem__day",children:r}),(0,s.jsx)("div",{className:"AgendaItem__time",children:l})]}),(0,s.jsxs)("div",{className:"AgendaItem__data",children:[(0,s.jsx)("h4",{children:a}),(0,s.jsxs)("p",{className:"p p-s",children:[" ",d]}),(0,s.jsx)("div",{className:c().description,children:(0,s.jsx)("p",{children:p})})]}),h&&(0,s.jsx)(o.default,{clickHandler:()=>_(n),children:"Archive"})]})},l=function(e){let[t,n]=(0,i.useState)(!1),[a,o]=(0,i.useState)(""),{day:l,highlights:d,admin:p}=e,[h,u]=(0,i.useState)(!1),f=[{id:"ABC-123",title:"Cached Entry",author:"Cached Author",description:"Cached Content"}],m=()=>{console.log("Querying /agenda/agenda-items/"),fetch("/api/agenda/agenda-items/").then(e=>e.json()).then(e=>{o(e),u(!1)}).catch(e=>{o(f),console.log(e)})},_=e=>{u(!0),n(!1),console.log("Archiving Agenda Item ..."+e),fetch("/api/agenda/agenda-items/"+e,{method:"DELETE",headers:{accept:"application/json"}}).then(e=>e.json()).then(()=>{m(),u(!1)}).catch(e=>{console.log(e),u(!1),n(!0)})};return(0,i.useEffect)(()=>{u(!0),m()},[o]),(0,s.jsx)("div",{children:(0,s.jsxs)("div",{className:"".concat(c().agendaList,"  ").concat(p?c().backoffice:" "," "),children:[a&&a.length>0&&a.map((e,t)=>(0,s.jsx)(r,{name:e.title,id:e.id,description:e.description,author:e.author,admin:p,handleArchive:_},t)),a&&0==a.length&&(0,s.jsx)("p",{children:"There are no confirmed talks just yet."})]})})}},52:function(e,t,n){"use strict";n.r(t);var s=n(7437),a=n(1693),c=n.n(a),i=n(1396),o=n.n(i);t.default=function(e){var t;let{children:n,link:a,external:i,inline:r,clickHandler:l,small:d,main:p,disabled:h,inverted:u}=e;return t=a?i?(0,s.jsxs)("a",{href:a,target:"_blank",children:["  ",(0,s.jsx)("span",{children:n})," "]}):(0,s.jsxs)(o(),{href:a,children:["  ",(0,s.jsx)("span",{children:n})," "]}):l?(0,s.jsxs)(o(),{className:"__container",href:"#",onClick:l,children:[" ",(0,s.jsx)("span",{children:n}),"  "]}):(0,s.jsxs)(o(),{href:"#",className:"__container",children:["  ",(0,s.jsx)("span",{children:n})," "]}),(0,s.jsx)("div",{className:c().button,children:t})}},1693:function(e){e.exports={button:"button_button__zBiTp"}},1449:function(e){e.exports={main:"agenda_main__TiUbp",hero:"agenda_hero__DClxi",agendaList:"agenda_agendaList__nj0AK",backoffice:"agenda_backoffice__Wlz5T",agendaItem:"agenda_agendaItem__1gGMa",open:"agenda_open__4fdYd",openTag:"agenda_openTag___xyig",description:"agenda_description__t0a8g"}},9617:function(e){e.exports={main:"backoffice_main__51xWt",hero:"backoffice_hero__Hwfjd",BackofficeContent:"backoffice_BackofficeContent__XwxMo",tabs:"backoffice_tabs__TzmKF"}},4902:function(e){e.exports={DebugList:"debug_DebugList__YiFgj",ServiceInfo:"debug_ServiceInfo__m5KK2",statusTag:"debug_statusTag__ulkHD",healthy:"debug_healthy__crQbC",unhealthy:"debug_unhealthy__Glymt",header:"debug_header__T5iow",description:"debug_description__JuDKU",descriptionItem:"debug_descriptionItem__8cqY0"}},6782:function(e){e.exports={EventsList:"events_EventsList__H1Y8V",EventItem:"events_EventItem__Jg_Cn",openTag:"events_openTag__tx1uO",header:"events_header__K3j9x",description:"events_description__IXpCv",open:"events_open__wCFok",codeContainer:"events_codeContainer__qSoh1"}},2880:function(e){e.exports={NotificationList:"notifications_NotificationList__j7E8c",NotificationItem:"notifications_NotificationItem__c2qkm",openTag:"notifications_openTag__S2TAb",header:"notifications_header__sFwWB",headerStatusTag:"notifications_headerStatusTag__XEkpW",approved:"notifications_approved__bFPWf",rejected:"notifications_rejected__A3PZL",description:"notifications_description__qSiGy",open:"notifications_open__dTwzw",descriptionTo:"notifications_descriptionTo__KhoKb",descriptionSubject:"notifications_descriptionSubject__YZwvH",descriptionBody:"notifications_descriptionBody__meYnW"}},7157:function(e){e.exports={main:"proposals_main__USy0x",hero:"proposals_hero__FIct1",ProposalList_Filters:"proposals_ProposalList_Filters__fv6zt",filterLabel:"proposals_filterLabel__i_cn3",filter:"proposals_filter__sP5Xo",inactive:"proposals_inactive__3_zuB",container:"proposals_container__IC4TY",ProposalList_Items:"proposals_ProposalList_Items__WFCAB",ProposalItem:"proposals_ProposalItem__SanIB",status:"proposals_status__lgADm",actions:"proposals_actions__5EjR_",description:"proposals_description__eQArS",pending:"proposals_pending__Q9fqt",archived:"proposals_archived__DyZj0",statusTag:"proposals_statusTag__smILf",approved:"proposals_approved__z7Yx5",rejected:"proposals_rejected__kkWbX"}}},function(e){e.O(0,[176,971,596,744],function(){return e(e.s=629)}),_N_E=e.O()}]);