console.log('kees client v0.0.1');

class KeesClient {
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
      let hook = this.setupHook(h);
      actions[hook.name] = hook.wrapper;
    });

    return actions;
  }

  setupHook(el){
    console.log(el)
    let name = el.attributes.hook.value;
    el.addEventListener('click', (el) => {
      this[name].bind(this)(el);
    });

    return {name: name, wrapper: this.elWrap(el)};
  }


  setupElements(){
    let elements = {}
    let check = document.querySelector(".getDevices");
    elements['getDevices'] = this.elWrap(check);

    let devices = document.querySelector(".devices");
    elements['devices'] = this.elWrap(devices);
    let deviceList = document.querySelector(".deviceList");
    elements['deviceList'] = this.elWrap(deviceList);

    let user = document.querySelector(".main .user .details");
    elements['user'] = this.elWrap(user);

    let deviceID = document.querySelector(".info .id");
    elements['deviceID'] = this.elWrap(deviceID);

    let logs = document.querySelector(".ws .logs");
    elements['logs'] = this.elWrap(logs);

    let actions = document.querySelector(".wsActions");
    elements['wsActions'] = this.elWrap(actions);

    return elements;
  }

  deviceSet(source){
    let fields = this.elements.device.el.querySelectorAll("input");
    console.log(fields);
  }


  restoreSession(){
    if (!localStorage.getItem('kees-client')){
      console.log('session not found, need to auth');
      return;
    }

    this.session = JSON.parse(localStorage.getItem('kees-client'));
    this.event("auth", "session", this.session);

    this.actions.auth.hide();
    this.actions.reset.show();

    this.elements.getDevices.show();
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


  getAuthInfo(){
    let el = this.elements.user.el;
    let payload = {
      username: el.querySelector('#username').value,
      password: el.querySelector('#password').value,
      client: el.querySelector('#client').value
    }

    return payload;
  }

  async auth(){
    let authInfo = this.getAuthInfo();
    this.event("auth", "auth info", authInfo);

    let resp = await fetch('/api/v1/auth', {
      method: 'POST',
      body: JSON.stringify(authInfo)
    });

    let data = await resp.json();
    this.event("auth", "auth response", data);

    localStorage.setItem('kees-client', JSON.stringify(data));
    this.restoreSession();
  }

  reset(){
    localStorage.removeItem("kees-client");
    this.event("auth", "reset", {});

    this.actions.reset.hide();
    this.elements.deviceID.clear();
    this.elements.check.hide();

    // TODO: clean up icon grabbing
    var icons = document.querySelectorAll(".info i").forEach(e => {
      e.classList.add("hide");
    });

    this.actions.auth.show();
  }

  async getDevices(){
    // TODO: clean up icon grabbing
    let status = this.elements.getDevices.el.querySelector("i");
    status.classList = 'fa-solid fa-circle-notch fa-spin';

    this.event("auth", "getDevices", this.session.data.jwt);

    let resp = await fetch('/api/v1/devices', {
      method: "GET",
      headers: {
        'X-Kees-JWT-Token': this.session.data.jwt.token
      }
    });

    let data = await resp.json();
    this.event("auth", "get devices response", data);
    status.classList = 'fa-solid fa-check good';


    this.renderDevices(data)
  }


  renderDevices(data){
    this.elements.devices.show();

    // for each device, render the device info and associated actions
    data.forEach(d => {
      console.log(d)
      var tmpl = `
        <div class="row device" data-device-id="${d.id}">
          <div class="row name">${d.name}</div>
          <div class="row id">${d.id}</div>
          <div class="row controller">${d.controller} @ ${d.version}</div>
          <div class="row actions">
          </div>
        </div>
      `;
      var t = document.createElement("template");
      t.innerHTML = tmpl;
      var controller = t.content.cloneNode(true);

      // grab actions from query selector because controller contents will be empty after we append below
      let actions = controller.querySelector(".actions");
      this.elements.deviceList.el.append(controller);

      // walk through the capabilities and creat+hook a button for them
      d.capabilities.forEach(c => {
        let btn = `
          <button data-device-id="${d.id}" data-device-action="${c}" hook="deviceAction">${c}</button>
        `;
        let t = document.createElement("template");
        t.innerHTML = btn;
        let action = t.content.cloneNode(true);
        actions.append(action);

        // ugly hack as document-fragment doesn't support attributes accessor
        action = actions.querySelector(`[data-device-action="${c}"]`);
        this.setupHook(action);
      })
    })
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
        token: this.session.data.jwt.token
      }
    }

    this.event("ws", "auth", auth);
    this.sendWS(auth)
  }


  deviceAction(e){
    let btn = e.srcElement;
    let action = btn.dataset.deviceAction;
    let deviceID = btn.dataset.deviceId;

    this.event("action", action, `${action} clicked for device:${deviceID}`);
    this.performAction(deviceID, action);
  }

  async performAction(deviceID, action){
    // TODO: clean up icon grabbing
    let status = this.elements.getDevices.el.querySelector("i");
    status.classList = 'fa-solid fa-circle-notch fa-spin';

    let endpoint = `/api/v1/devices/${deviceID}/commands/${action}`;
    this.event("action", action, {token: this.session.data.jwt, endpoint: endpoint});

    let resp = await fetch(endpoint, {
      method: "POST",
      headers: {
        'X-Kees-JWT-Token': this.session.data.jwt.token
      }
    });

    let data = await resp.json();
    this.event("auth", `${action} on ${deviceID} response`, data);
    status.classList = 'fa-solid fa-check good';
  }
}

var k = new KeesClient
