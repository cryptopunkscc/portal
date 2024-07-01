var Lt=Object.defineProperty;var It=(n,t,e)=>t in n?Lt(n,t,{enumerable:!0,configurable:!0,writable:!0,value:e}):n[t]=e;var M=(n,t,e)=>(It(n,typeof t!="symbol"?t+"":t,e),e),Jt=(n,t,e)=>{if(!t.has(n))throw TypeError("Cannot "+e)};var G=(n,t,e)=>{if(t.has(n))throw TypeError("Cannot add the same private member more than once");t instanceof WeakSet?t.add(n):t.set(n,e)};var P=(n,t,e)=>(Jt(n,t,"access private method"),e);(function(){const t=document.createElement("link").relList;if(t&&t.supports&&t.supports("modulepreload"))return;for(const s of document.querySelectorAll('link[rel="modulepreload"]'))r(s);new MutationObserver(s=>{for(const o of s)if(o.type==="childList")for(const a of o.addedNodes)a.tagName==="LINK"&&a.rel==="modulepreload"&&r(a)}).observe(document,{childList:!0,subtree:!0});function e(s){const o={};return s.integrity&&(o.integrity=s.integrity),s.referrerpolicy&&(o.referrerPolicy=s.referrerpolicy),s.crossorigin==="use-credentials"?o.credentials="include":s.crossorigin==="anonymous"?o.credentials="omit":o.credentials="same-origin",o}function r(s){if(s.ep)return;s.ep=!0;const o=e(s);fetch(s.href,o)}})();function g(){}function $t(n){return n()}function ht(){return Object.create(null)}function I(n){n.forEach($t)}function bt(n){return typeof n=="function"}function Y(n,t){return n!=n?t==t:n!==t||n&&typeof n=="object"||typeof n=="function"}let B;function dt(n,t){return B||(B=document.createElement("a")),B.href=t,n===B.href}function kt(n){return Object.keys(n).length===0}function Mt(n,...t){if(n==null)return g;const e=n.subscribe(...t);return e.unsubscribe?()=>e.unsubscribe():e}function Pt(n,t,e){n.$$.on_destroy.push(Mt(t,e))}function h(n,t){n.appendChild(t)}function st(n,t,e){n.insertBefore(t,e||null)}function Z(n){n.parentNode&&n.parentNode.removeChild(n)}function Bt(n,t){for(let e=0;e<n.length;e+=1)n[e]&&n[e].d(t)}function _(n){return document.createElement(n)}function H(n){return document.createTextNode(n)}function U(){return H(" ")}function Ht(n,t,e,r){return n.addEventListener(t,e,r),()=>n.removeEventListener(t,e,r)}function p(n,t,e){e==null?n.removeAttribute(t):n.getAttribute(t)!==e&&n.setAttribute(t,e)}function Qt(n){return Array.from(n.childNodes)}function V(n,t){t=""+t,n.data!==t&&(n.data=t)}let x;function E(n){x=n}function Wt(){if(!x)throw new Error("Function called outside component initialization");return x}function At(n){Wt().$$.on_destroy.push(n)}const O=[],_t=[];let N=[];const gt=[],zt=Promise.resolve();let tt=!1;function Ft(){tt||(tt=!0,zt.then(qt))}function et(n){N.push(n)}const X=new Set;let A=0;function qt(){if(A!==0)return;const n=x;do{try{for(;A<O.length;){const t=O[A];A++,E(t),Yt(t.$$)}}catch(t){throw O.length=0,A=0,t}for(E(null),O.length=0,A=0;_t.length;)_t.pop()();for(let t=0;t<N.length;t+=1){const e=N[t];X.has(e)||(X.add(e),e())}N.length=0}while(O.length);for(;gt.length;)gt.pop()();tt=!1,X.clear(),E(n)}function Yt(n){if(n.fragment!==null){n.update(),I(n.before_update);const t=n.dirty;n.dirty=[-1],n.fragment&&n.fragment.p(n.ctx,t),n.after_update.forEach(et)}}function Zt(n){const t=[],e=[];N.forEach(r=>n.indexOf(r)===-1?t.push(r):e.push(r)),e.forEach(r=>r()),N=t}const Q=new Set;let b;function Dt(){b={r:0,c:[],p:b}}function Kt(){b.r||I(b.c),b=b.p}function j(n,t){n&&n.i&&(Q.delete(n),n.i(t))}function W(n,t,e,r){if(n&&n.o){if(Q.has(n))return;Q.add(n),b.c.push(()=>{Q.delete(n),r&&(e&&n.d(1),r())}),n.o(t)}else r&&r()}function Ot(n){n&&n.c()}function ot(n,t,e,r){const{fragment:s,after_update:o}=n.$$;s&&s.m(t,e),r||et(()=>{const a=n.$$.on_mount.map($t).filter(bt);n.$$.on_destroy?n.$$.on_destroy.push(...a):I(a),n.$$.on_mount=[]}),o.forEach(et)}function at(n,t){const e=n.$$;e.fragment!==null&&(Zt(e.after_update),I(e.on_destroy),e.fragment&&e.fragment.d(t),e.on_destroy=e.fragment=null,e.ctx=[])}function Tt(n,t){n.$$.dirty[0]===-1&&(O.push(n),Ft(),n.$$.dirty.fill(0)),n.$$.dirty[t/31|0]|=1<<t%31}function it(n,t,e,r,s,o,a,l=[-1]){const i=x;E(n);const c=n.$$={fragment:null,ctx:[],props:o,update:g,not_equal:s,bound:ht(),on_mount:[],on_destroy:[],on_disconnect:[],before_update:[],after_update:[],context:new Map(t.context||(i?i.$$.context:[])),callbacks:ht(),dirty:l,skip_bound:!1,root:t.target||i.$$.root};a&&a(c.root);let y=!1;if(c.ctx=e?e(n,t.props||{},(d,S,...w)=>{const v=w.length?w[0]:S;return c.ctx&&s(c.ctx[d],c.ctx[d]=v)&&(!c.skip_bound&&c.bound[d]&&c.bound[d](v),y&&Tt(n,d)),S}):[],c.update(),y=!0,I(c.before_update),c.fragment=r?r(c.ctx):!1,t.target){if(t.hydrate){const d=Qt(t.target);c.fragment&&c.fragment.l(d),d.forEach(Z)}else c.fragment&&c.fragment.c();t.intro&&j(n.$$.fragment),ot(n,t.target,t.anchor,t.customElement),qt()}E(i)}class ct{$destroy(){at(this,1),this.$destroy=g}$on(t,e){if(!bt(e))return g;const r=this.$$.callbacks[t]||(this.$$.callbacks[t]=[]);return r.push(e),()=>{const s=r.indexOf(e);s!==-1&&r.splice(s,1)}}$set(t){this.$$set&&!kt(t)&&(this.$$.skip_bound=!0,this.$$set(t),this.$$.skip_bound=!1)}}const u={};function lt(n,t){n!==void 0&&Object.assign(u,{platform:n,...t()})}let f;try{f=window}catch{f={}}const Gt=typeof f.go>"u"?void 0:"wails",Ut=()=>({astral_conn_accept:f.go.main.Adapter.ConnAccept,astral_conn_close:f.go.main.Adapter.ConnClose,astral_conn_read:f.go.main.Adapter.ConnRead,astral_conn_write:f.go.main.Adapter.ConnWrite,astral_node_info:f.go.main.Adapter.NodeInfo,astral_query:f.go.main.Adapter.Query,astral_query_name:f.go.main.Adapter.QueryName,astral_resolve:f.go.main.Adapter.Resolve,astral_service_close:f.go.main.Adapter.ServiceClose,astral_service_register:f.go.main.Adapter.ServiceRegister,astral_interrupt:f.go.main.Adapter.Interrupt,sleep:f.go.main.Adapter.Sleep,log:async(...n)=>await f.go.main.Adapter.LogArr(n)});lt(Gt,Ut);const Vt=typeof _app_host>"u"?void 0:"android",Xt=()=>{const n=new Map;window._resolve=(e,r)=>{n.get(e)[0](r),n.delete(e)},window._reject=(e,r)=>{n.get(e)[1](r),n.delete(e)};const t=e=>new Promise((r,s)=>n.set(e(),[r,s]));return{astral_node_info:e=>t(()=>_app_host.nodeInfo(e)).then(r=>JSON.parse(r)),astral_conn_accept:e=>t(()=>_app_host.connAccept(e)),astral_conn_close:e=>t(()=>_app_host.connClose(e)),astral_conn_read:e=>t(()=>_app_host.connRead(e)),astral_conn_write:(e,r)=>t(()=>_app_host.connWrite(e,r)),astral_query:(e,r)=>t(()=>_app_host.query(e,r)),astral_query_name:(e,r)=>t(()=>_app_host.queryName(e,r)),astral_resolve:e=>t(()=>_app_host.resolve(e)),astral_service_close:e=>t(()=>_app_host.serviceClose(e)),astral_service_register:e=>t(()=>_app_host.serviceRegister(e)),astral_interrupt:()=>t(()=>_app_host.interrupt()),sleep:e=>t(()=>_app_host.sleep(e)),log:e=>_app_host.logArr(JSON.stringify(e))}};lt(Vt,Xt);const te=typeof _log>"u"?void 0:"common",ee=()=>({astral_conn_accept:_astral_conn_accept,astral_conn_close:_astral_conn_close,astral_conn_read:_astral_conn_read,astral_conn_write:_astral_conn_write,astral_node_info:_astral_node_info,astral_query:_astral_query,astral_query_name:_astral_query_name,astral_resolve:_astral_resolve,astral_service_close:_astral_service_close,astral_service_register:_astral_service_register,astral_interrupt:_astral_interrupt,sleep:_sleep,log:_log});lt(te,ee);class Nt{async register(t){return await u.astral_service_register(t),new ne(t)}async query(t,e){e=e||"";const r=await u.astral_query(e,t),s=JSON.parse(r);return new z(s,t)}async queryName(t,e){const r=await u.astral_query_name(t,e),s=JSON.parse(r);return new z(s,e)}async nodeInfo(t){return await u.astral_node_info(t)}async resolve(t){return await u.astral_resolve(t)}async interrupt(){await u.astral_interrupt()}}class ne{constructor(t){this.port=t}async accept(){const t=await u.astral_conn_accept(this.port),e=JSON.parse(t);return new z(e)}async close(){await u.astral_service_close(this.port)}}class z{constructor(t){this.id=t.id,this.query=t.query,this.remoteId=t.remoteId}async read(){try{return await u.astral_conn_read(this.id)}catch(t){throw this.done=!0,t}}async write(t){try{return await u.astral_conn_write(this.id,t)}catch(e){throw this.done=!0,e}}async close(){this.done=!0,await u.astral_conn_close(this.id)}}function jt(n,t){const e=re(t),r=n.copy();for(let[s,o]of e){if(r[s])throw`method '${s}' already exist`;r[s]=r.call(o)}return r}const mt=/^\*/;function re(n){if(!Array.isArray(n))throw`cannot prepare routes of type ${typeof n}`;const t=[];for(let e in n){const r=n[e];switch(typeof r){case"string":const s=r.replace(mt,"");t.push([s,s]);continue;case"object":for(let o in r)for(let a of r[o]){a=a.replace(mt,"");const l=[o,a].join(".");t.push([a,l])}}}return t}function St(n,t,...e){const r=new se(n,t,e);return Object.assign(async(...o)=>await r.request(...o),{inner:r,map:(...o)=>r.map(...o),filter:(...o)=>r.filter(...o),request:async(...o)=>await r.request(...o),collect:async(...o)=>await r.collect(...o),conn:async(...o)=>await r.conn(...o)})}var L,nt;class se{constructor(t,e,r){G(this,L);M(this,"mapper",t=>t);M(this,"params",[]);M(this,"single",!0);this.client=t,this.port=e,this.params=Array.isArray(r)?[]:r}map(t){const e=this.mapper;return this.mapper=r=>t(e(r)),this}filter(t){return this.map(e=>{if(t(e))return e})}async request(...t){return t.length>0&&(this.params=t),await P(this,L,nt).call(this,async e=>await e.request(...t))}async collect(...t){return t.length>0&&(this.params=t),await P(this,L,nt).call(this,async e=>await e.collect(...t))}async conn(...t){const e=t.length>0?t:this.params;return this.client.conn(this.port,...e)}}L=new WeakSet,nt=async function(t){const e=await this.conn();return e.mapper=this.mapper,this.result=await t(e),this.mapper=r=>r,this.single&&await e.close().catch(u.log),this.result};function oe(n){const t=n.search(/[?.{\[]/);if(t===-1)return[n];const e=n.slice(0,t);let r=n.slice(t,n.length);return/^[.?]/.test(r)&&(r=r.slice(1)),[e,r]}function ae(n){return n.search(/[?{\[]/)>-1}var F,Ct;class Et extends z{constructor(e){super(e);G(this,F)}bind(...e){return jt(this,e)}copy(){return this}call(e,...r){const s=St(this,P(this,F,Ct).call(this,e),...r);return s.inner.single=!1,s}map(e){if(this.mapper){const r=this.mapper;this.mapper=s=>e(r(s))}else this.mapper=e;return this}async conn(e,...r){let s=e||"";return r.length>0&&(s&&(s+="?"),s+=JSON.stringify(r)),s&&await this.write(s+`
`),this}async encode(e){let r=JSON.stringify(e);return r===void 0&&(r="{}"),await super.write(r+`
`)}async decode(){const e=await this.read(),r=JSON.parse(e);if(r===null)return null;if(r.error)throw r.error;return r}async request(...e){const r=this.mapper;for(this.result=null;;){const s=await this.decode();if(s===void 0)continue;if(s===null)return this.result;if(this.result=s,!r)return s;const o=await r(s);if(o!==void 0)return o===null?this.result:o}}async collect(...e){const r=this.mapper?this.mapper:null;let s;for(r?s=async o=>{if(o=await r.call(this,o),o===null)return this.result;o&&this.result.push(o)}:s=o=>this.result.push(o),this.result=[];;){let o=await this.decode();if(o===null)return this.result;s(o)}}}F=new WeakSet,Ct=function(e){if(ae(this.query))throw`cannot nest connection for complete query ${chunks}`;return e};function ie(n){let t=Rt(n.handlers);return t=ce(t),t=le(t,n.routes),t}function Rt(n,...t){if(typeof n!="object")return t;const e=Object.getOwnPropertyNames(n);if(e.length===0)return t;const r=[];for(let s of e){const o=n[s],a=Rt(o,...t,s);typeof a[0]=="string"?r.push(a):r.push(...a)}return r}function ce(n){const t=[];for(let e of n)t.push(e.join("."));return t}function le(n,t){t=t||[];let e=[...n];for(let r of t){if(r==="*")return[t];const s=r.length-1;/[*:]/.test(r.slice(s))&&(r=r.slice(0,s)),e=e.filter(o=>!o.startsWith(r))}return t=t.filter(r=>!r.endsWith(":")),e.push(...t),e}async function ue(n,t){const e=ie(t);for(let r of e){const s=await n.register(r);fe(t,s).catch(u.log)}}async function fe(n,t){for(;;){let e=await t.accept();e=new Et(e),pe(n,e).catch(u.log).finally(()=>e.close().catch(u.log))}}async function pe(n,t){const e={...n.handlers,...n.inject,conn:t},r=t.query;let[s,o]=rt(n.handlers,r),a=s,l,i;for(;;){if(i=typeof a=="function",o&&!i){await t.encode({error:`no handler for query ${o} ${typeof a}`});return}if(o||i){try{l=await he(e,a,o)}catch(c){l={error:c}}await t.encode(l),a=s}o=await t.read(),typeof a=="object"&&([a,o]=rt(a,o))}}async function he(n,t,e){const r=typeof t;switch(r){case"function":if(!e)return await t(n);const s=JSON.parse(e);return Array.isArray(s)?await t(...s,n):await t(s,n);case"object":return;default:throw`invalid handler type ${r}`}}function rt(n,t){if(t==="")return[n];const[e,r]=oe(t),s=n[e];if(r===void 0)return[s];if(typeof s<"u")return rt(s,r);if(typeof n=="function")return[n,r];throw"cannot unfold"}class ut extends Nt{bind(...t){return jt(this,t)}copy(t){return Object.assign(new ut,{...this,...t})}target(t){return this.targetId=t,this}call(t,...e){return St(this,t,...e)}async conn(t,...e){const r=de(t,e),s=await super.query(r,this.targetId);return new Et(s)}async serve(t){return await ue(this,t)}}function de(n,t){let e=n;return t.length>0&&(e+="?"+JSON.stringify(t)),e}const{log:C,sleep:Ce,platform:Re}=u,_e=new Nt,xt=new ut;C("launcher start");const R=xt.bind({portal:["open","install","uninstall"]});R.observe=async()=>{const n=await xt.conn("portal.observe");return{next:async()=>await n.decode(),more:async t=>await n.encode(t),close:n.close}};class ge{constructor(){this.launch=R.open,this.install=R.install,this.uninstall=R.uninstall}}function me(n){let t,e,r,s,o,a,l,i,c,y=n[0].title+"",d,S,w,v=n[0].installed?"(Zainstalowano)":"",D,ft,J,k=n[0].description+"",K,T,pt;return{c(){t=_("div"),e=_("div"),r=_("img"),a=U(),l=_("div"),i=_("div"),c=_("h2"),d=H(y),S=U(),w=_("span"),D=H(v),ft=U(),J=_("p"),K=H(k),dt(r.src,s=n[0].icon)||p(r,"src",s),p(r,"alt",o=n[0].name),p(r,"class","svelte-1vqrebr"),p(e,"class","app-icon svelte-1vqrebr"),p(c,"class","svelte-1vqrebr"),p(w,"class","installed-status svelte-1vqrebr"),p(i,"class","app-name svelte-1vqrebr"),p(J,"class","app-description svelte-1vqrebr"),p(l,"class","app-info svelte-1vqrebr"),p(t,"class","app-item svelte-1vqrebr")},m(m,$){st(m,t,$),h(t,e),h(e,r),h(t,a),h(t,l),h(l,i),h(i,c),h(c,d),h(i,S),h(i,w),h(w,D),h(l,ft),h(l,J),h(J,K),T||(pt=Ht(t,"click",n[1]),T=!0)},p(m,[$]){$&1&&!dt(r.src,s=m[0].icon)&&p(r,"src",s),$&1&&o!==(o=m[0].name)&&p(r,"alt",o),$&1&&y!==(y=m[0].title+"")&&V(d,y),$&1&&v!==(v=m[0].installed?"(Zainstalowano)":"")&&V(D,v),$&1&&k!==(k=m[0].description+"")&&V(K,k)},i:g,o:g,d(m){m&&Z(t),T=!1,pt()}}}function ye(n,t,e){const r=new ge;let{app:s={title:"",package:"",name:"",description:"",icon:"",url:"",installed:!1}}=t;const o=()=>r.launch(s.package).catch(C);return n.$$set=a=>{"app"in a&&e(0,s=a.app)},[s,o]}class we extends ct{constructor(t){super(),it(this,t,ye,me,Y,{app:0})}}const q=[];function ve(n,t=g){let e;const r=new Set;function s(l){if(Y(n,l)&&(n=l,e)){const i=!q.length;for(const c of r)c[1](),q.push(c,n);if(i){for(let c=0;c<q.length;c+=2)q[c][0](q[c+1]);q.length=0}}}function o(l){s(l(n))}function a(l,i=g){const c=[l,i];return r.add(c),r.size===1&&(e=t(s)||g),l(n),()=>{r.delete(c),r.size===0&&e&&(e(),e=null)}}return{set:s,update:o,subscribe:a}}class $e{constructor(){this.apps=[],this.store=ve([]),this.channel=null,R.observe().catch(console.log).then(t=>{this.channel=t,this.run().catch(console.log),this.loadMore()}),console.log(this)}async run(){for(;this.channel;){let t;try{t=await this.channel.next()}catch(e){throw C("error: ",JSON.stringify(e)),e}this.apps.push(t),this.store.set(this.apps)}C("close run")}subscribe(t,e){return this.store.subscribe(t,e)}loadMore(t){var e,r;t=t||10,(r=(e=this.channel)==null?void 0:e.more(t))==null||r.catch(console.log)}cancel(){var t,e;C("cancel"),(e=(t=this.channel)==null?void 0:t.close())==null||e.catch(console.log),this.channel=null}}function yt(n){n=n||1;const t=window.scrollY,e=document.documentElement.scrollHeight,r=window.innerHeight;return e-t-r<n}function be(n){let t=yt(),e=t;e&&n();const r=()=>{e=yt(),!t&&e&&n(),t=e};return document.addEventListener("scroll",r),()=>document.removeEventListener("scroll",r)}function wt(n,t,e){const r=n.slice();return r[2]=t[e],r[4]=e,r}function vt(n){let t,e;return t=new we({props:{app:n[2]}}),{c(){Ot(t.$$.fragment)},m(r,s){ot(t,r,s),e=!0},p(r,s){const o={};s&1&&(o.app=r[2]),t.$set(o)},i(r){e||(j(t.$$.fragment,r),e=!0)},o(r){W(t.$$.fragment,r),e=!1},d(r){at(t,r)}}}function Ae(n){let t,e,r=n[0],s=[];for(let a=0;a<r.length;a+=1)s[a]=vt(wt(n,r,a));const o=a=>W(s[a],1,1,()=>{s[a]=null});return{c(){t=_("div");for(let a=0;a<s.length;a+=1)s[a].c();p(t,"class","all apps svelte-gkucmm")},m(a,l){st(a,t,l);for(let i=0;i<s.length;i+=1)s[i]&&s[i].m(t,null);e=!0},p(a,[l]){if(l&1){r=a[0];let i;for(i=0;i<r.length;i+=1){const c=wt(a,r,i);s[i]?(s[i].p(c,l),j(s[i],1)):(s[i]=vt(c),s[i].c(),j(s[i],1),s[i].m(t,null))}for(Dt(),i=r.length;i<s.length;i+=1)o(i);Kt()}},i(a){if(!e){for(let l=0;l<r.length;l+=1)j(s[l]);e=!0}},o(a){s=s.filter(Boolean);for(let l=0;l<s.length;l+=1)W(s[l]);e=!1},d(a){a&&Z(t),Bt(s,a)}}}function qe(n,t,e){let r;const s=new $e;return Pt(n,s,o=>e(0,r=o)),be(()=>s.loadMore()),At(()=>s.cancel()),[r,s]}class Oe extends ct{constructor(t){super(),it(this,t,qe,Ae,Y,{})}}function Ne(n){let t,e,r;return e=new Oe({}),{c(){t=_("main"),Ot(e.$$.fragment),p(t,"class","svelte-vx7f26")},m(s,o){st(s,t,o),ot(e,t,null),r=!0},p:g,i(s){r||(j(e.$$.fragment,s),r=!0)},o(s){W(e.$$.fragment,s),r=!1},d(s){s&&Z(t),at(e)}}}function je(n){return At(_e.interrupt),[]}class Se extends ct{constructor(t){super(),it(this,t,je,Ne,Y,{})}}new Se({target:document.getElementById("app")});