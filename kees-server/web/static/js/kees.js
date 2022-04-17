console.log('kees websocket test client v0.0.1');

var session = null;
var ws = null;
var logCount = 0;

var actions = setupHooks();
restoreSession();


function setupHooks(){
  let hooks = document.querySelectorAll('[hook]');
  let actions = {};
  hooks.forEach(h => {
    let name = h.attributes.hook.value;
    actions[name] = {
      el: h,
      toggle: () => { toggleElement(h); },
      show:   () => { showElement(h); },
      hide:   () => { hideElement(h); }
    };

    h.addEventListener('click', window[name]);
  });

  return actions;
}

function toggleElement(el){ el.classList.toggle('hide'); }
function showElement(el)  { el.classList.remove('hide'); }
function hideElement(el)  { el.classList.add('hide'); }


function auth(){
  console.log("Authing Device");
  var details = getDetails();
  console.log(details);

  submitAuth(details.token, details.device);
}

function reset(){
  console.log("clearing session");
  localStorage.removeItem("kees");

  addWSLog("auth", "reset", {});

  if (ws){
    ws.close();
  }

  actions["reset"].hide();

  var id = document.querySelector(".info .id").innerHTML = '';
  var check = document.querySelector(".check");
  check.classList.add("hide");

  var wsStart = document.querySelector(".wsStart");
  wsStart.classList.add("hide");

  var icons = document.querySelectorAll(".info i").forEach(e => {
    e.classList.add("hide");
  });

  actions["auth"].show();
}

async function check(){
  var status = document.querySelector(".check i");
  status.classList = 'fa-solid fa-circle-notch fa-spin';

  addWSLog("auth", "check", session.jwt);

  let resp = await fetch('/ws/v1/auth/check', {
    method: "GET",
    headers: {
      'X-Kees-JWT-Token': session.jwt.token
    }
  });

  let data = await resp.json();
  console.log(data);
  addWSLog("auth", "check ack", data);
  status.classList = 'fa-solid fa-check good';
}


function getDetails(){
  console.log("Getting Device Details");
  var details = document.querySelector(".main .device .details");
  var deviceInfo = {
    name: details.querySelector("#name").value,
    version: details.querySelector("#version").value,
    controller: details.querySelector("#controller").value
  }

  var token = details.querySelector("#deviceToken").value;

  var payload = {
    device: deviceInfo,
    token: token
  }
  return payload;
}

async function submitAuth(token, payload){
  addWSLog("auth", "auth", {token: token, payload: payload});

  let resp = await fetch('/ws/v1/auth', {
    method: 'POST',
    headers: {
      'X-Kees-MC-Token': token
    },
    body: JSON.stringify(payload)
  });

  console.log(resp);

  let data = await resp.json();
  console.log(data);

  addWSLog("auth", "auth ack", data);

  localStorage.setItem('kees', JSON.stringify(data));
  restoreSession();
}


function restoreSession(){
  if (!localStorage.getItem('kees')){
    console.log('session not found, need to auth');
    return;
  }

  console.log('session found, restoring');
  session = JSON.parse(localStorage.getItem('kees'));
  console.log(session);

  addWSLog("auth", "session", session);

  actions["auth"].hide();

  actions["reset"].show();


  var deviceInfo = document.querySelector(".main .device .details");
  deviceInfo.querySelector("#name").value = session.device.name;
  deviceInfo.querySelector("#version").value = session.device.version;
  deviceInfo.querySelector("#controller").value = session.device.controller;

  var id = document.querySelector(".info .id").innerHTML = `
    <strong>device id</strong><br> ${session.device.id}
  `;

  var check = document.querySelector(".check");
  check.classList.remove("hide");

  var wsStart = document.querySelector(".wsStart");
  wsStart.classList.remove("hide");
}

function addWSLog(category, type, data){
  var logs = document.querySelector(".ws .logs");

  logCount++;
  var tmpl = logTemplate({
    category: category,
    type: type,
    data: data
  });

  logs.prepend(tmpl);
}


function logTemplate(data){
  var gray = (logCount % 2) == 0 ? 'gray' : '';
  var tmpl = `
    <div class="row event ${gray}">
      <div class="four columns">
        <div class="row category">${data.category}</div>
        <div class="row type">${data.type}</div>
        <div class="row date">${(new Date()).toLocaleString()}</div>
      </div>
      <div class="eight columns data">
        <pre>${JSON.stringify(data.data, undefined, 2)}</pre>
      </div>
    </div>
  `;

  var t = document.createElement("template");
  t.innerHTML = tmpl;
  var el = t.content.cloneNode(true);
  return el;
}

function startWS(){
  var status = document.querySelector(".wsStart i");
  status.classList = 'fa-solid fa-circle-notch fa-spin';

  var host = `ws://${document.location.host}/ws/v1/mc`;
  console.log(`connecting to ${host}`);
  ws = new WebSocket(host);
  console.log(ws);

  ws.onmessage = wsMessage;
  ws.onerror = wsError;
  ws.onclose = wsClose;
  ws.onopen = wsOpen;
}

function wsOpen(e){
  console.log("ws open");
  console.log(e);

  var status = document.querySelector(".wsStart i");
  status.classList = 'fa-solid fa-check good';

  var wsActions = document.querySelector(".wsAction");
  wsActions.classList.remove("hide");

  var ws = document.querySelector(".ws");
  ws.classList.remove("hide");


  var payload = {url: e.target.url};
  addWSLog("ws", "open", payload);
}

function wsMessage(e){
  console.log("ws message");
  console.log(e);
  var payload = JSON.parse(e.data);
  addWSLog("ws", "inbound", payload);
}

function wsError(e){
  console.log("ws error");
  console.log(e);
  addWSLog("ws", "error", e.data)
}

function wsClose(e){
  console.log("ws close");
  console.log(e);
  var payload = {reason: e.reason, code: e.code, timestamp: e.timestamp, type: e.type, clean: e.wasClean}
  addWSLog("ws", "close", payload);

  ws = null;
  var status = document.querySelector(".wsStart i");
  status.classList = 'fa-solid fa-times bad';

  var wsActions = document.querySelector(".wsAction");
  wsActions.classList.add("hide");

}

function badAuth(){
  badauth = {
    message: "auth",
    token: "eatmyshorts"
  }
  addWSLog("ws", "auth", badauth);
  sendWS(badauth);
}


function goodAuth(){
  auth = {
    message: "auth",
    data: {
      token: session.jwt.token
    }
  }

  addWSLog("ws", "auth", auth);
  sendWS(auth)
}

function sendWS(data){
  var payload = JSON.stringify(data);
  ws.send(payload);
}
