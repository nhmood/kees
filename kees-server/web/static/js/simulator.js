console.log('kees websocket client simulator v0.0.1');

class KeesClientSimulator {
  constructor(){
    this.session = undefined;
    this.ws = undefined;
    this.logCount = 0;

    this.actions = this.setupHooks();
    this.elements = this.setupElements();
    this.restoreSession();
  }

  elWrap(el){
    return {
      el: el,
      toggle: () => { this.toggleElement(el); },
      show:   () => { this.showElement(el); },
      hide:   () => { this.hideElement(el); },
      clear:  () => { this.clearElement(el); }
    }
  }

  toggleElement(el){ el.classList.toggle('hide'); }
  showElement(el)  { el.classList.remove('hide'); }
  hideElement(el)  { el.classList.add('hide'); }
  clearElement(el) { el.innerHTML = ''; }

  setupHooks(){
    let hooks = document.querySelectorAll('[hook]');
    let actions = {};
    hooks.forEach(h => {
      let name = h.attributes.hook.value;
      actions[name] = this.elWrap(h);
      h.addEventListener('click', (e) => {
        this[name].bind(this)(e);
      });
    });

    return actions;
  }


  setupElements(){
    let elements = {}
    let check = document.querySelector(".check");
    elements['check'] = this.elWrap(check);

    let ws = document.querySelector(".startWS");
    elements['startWS'] = this.elWrap(ws);

    let device = document.querySelector(".main .device .details");
    elements['device'] = this.elWrap(device);

    let deviceID = document.querySelector(".info .id");
    elements['deviceID'] = this.elWrap(deviceID);

    let logs = document.querySelector(".ws .logs");
    elements['logs'] = this.elWrap(logs);

    let actions = document.querySelector(".wsActions");
    elements['wsActions'] = this.elWrap(actions);

    return elements;
  }

  setupDevice(device){
    this.device = {
      id: device.id,
      name: device.name,
      version: device.version,
      controller: device.controller
    };


    return device;
  }

  setupDeviceUI(){
    let el = this.elements.device.el;
    el.querySelector('#name').value = this.device.name;
    el.querySelector('#version').value = this.device.version;
    el.querySelector('#controller').value = this.device.controller;

    this.elements.deviceID.el.innerHTML = `
      <strong>device id</strong><br> ${this.device.id}
    `;
  }

  deviceSet(source){
    let fields = this.elements.device.el.querySelectorAll("input");
    console.log(fields);
  }


  restoreSession(){
    if (!localStorage.getItem('kees')){
      console.log('session not found, need to auth');
      return;
    }

    this.session = JSON.parse(localStorage.getItem('kees'));
    this.event("auth", "session", this.session);

    this.actions.auth.hide();
    this.actions.reset.show();

    this.device = this.setupDevice(this.session.device);
    this.setupDeviceUI();

    this.elements.check.show();
    this.elements.startWS.show();
  }


  event(category, type, data){
    this.logCount++;
    var gray = (this.logCount % 2) == 0 ? 'gray' : '';
    var tmpl = `
      <div class="row event ${gray}">
        <div class="four columns">
          <div class="row category">${category}</div>
          <div class="row type">${type}</div>
          <div class="row date">${(new Date()).toLocaleString()}</div>
        </div>
        <div class="eight columns data">
          <pre>${JSON.stringify(data, undefined, 2)}</pre>
        </div>
      </div>
    `;

    var t = document.createElement("template");
    t.innerHTML = tmpl;
    var el = t.content.cloneNode(true);
    this.elements.logs.el.prepend(el);
  }


  getDeviceInfo(){
    let el = this.elements.device.el;
    let device = {
      name: el.querySelector('#name').value,
      version: el.querySelector('#version').value,
      controller: el.querySelector('#controller').value,
      capabilities: [
        "play",
        "stop",
        "fast_foward",
        "rewind",
        "pause"
      ]
    }

    let token = el.querySelector("#deviceToken").value;
    let payload = {
      device: device,
      token: token
    };

    return payload;
  }

  async auth(){
    let deviceInfo = this.getDeviceInfo();
    this.event("auth", "device info", deviceInfo);

    let resp = await fetch('/ws/v1/auth', {
      method: 'POST',
      headers: {
        'X-Kees-MC-Token': deviceInfo.token
      },
      body: JSON.stringify(deviceInfo.device)
    });

    let data = await resp.json();
    this.event("auth", "auth ack", data);

    localStorage.setItem('kees', JSON.stringify(data));
    this.restoreSession();
  }

  reset(){
    localStorage.removeItem("kees");
    this.event("auth", "reset", {});
    if (this.ws){ this.ws.close() };

    this.actions.reset.hide();
    this.elements.deviceID.clear();
    this.elements.check.hide();
    this.elements.startWS.hide();

    // TODO: clean up icon grabbing
    var icons = document.querySelectorAll(".info i").forEach(e => {
      e.classList.add("hide");
    });

    this.actions.auth.show();
  }

  async check(){
    // TODO: clean up icon grabbing
    let status = this.elements.check.el.querySelector("i");
    status.classList = 'fa-solid fa-circle-notch fa-spin';

    this.event("auth", "check", this.session.jwt);

    let resp = await fetch('/ws/v1/auth/check', {
      method: "GET",
      headers: {
        'X-Kees-JWT-Token': this.session.jwt.token
      }
    });

    let data = await resp.json();
    this.event("auth", "check ack", data);
    status.classList = 'fa-solid fa-check good';
  }


  startWS(){
    var status = this.elements.startWS.el.querySelector("i");
    status.classList = 'fa-solid fa-circle-notch fa-spin';

    var host = `ws://${document.location.host}/ws/v1/mc`;
    this.ws = new WebSocket(host);

    this.ws.onmessage = this.wsMessage.bind(this);
    this.ws.onerror   = this.wsError.bind(this);
    this.ws.onclose   = this.wsClose.bind(this);
    this.ws.onopen    = this.wsOpen.bind(this);
  }

  sendWS(data){
    var payload = JSON.stringify(data);
    this.ws.send(payload);
  }

  wsOpen(e){
    var status = this.elements.startWS.el.querySelector("i");
    status.classList = 'fa-solid fa-check good';

    this.elements.wsActions.show();
    this.event("ws", "open", {url: e.target.url});
  }

  wsMessage(e){
    var payload = JSON.parse(e.data);
    this.event("ws", "inbound", payload);
  }

  wsError(e){
    this.event("ws", "error", e.data)
  }

  wsClose(e){
    var payload = {reason: e.reason, code: e.code, timestamp: e.timestamp, type: e.type, clean: e.wasClean}
    this.event("ws", "close", payload);

    this.ws = null;
    var status = this.elements.startWS.el.querySelector("i");
    status.classList = 'fa-solid fa-times bad';

    this.elements.wsActions.hide();
  }

  badAuth(){
    let badauth = {
      state: "auth",
      message: "this is a bad auth",
      token: "eatmyshorts"
    }
    this.event("ws", "auth", badauth);
    this.sendWS(badauth);
  }


  goodAuth(){
    let auth = {
      state: "auth",
      message: "this is a good auth",
      data: {
        token: this.session.jwt.token
      }
    }

    this.event("ws", "auth", auth);
    this.sendWS(auth)
  }


}

var k = new KeesClientSimulator
