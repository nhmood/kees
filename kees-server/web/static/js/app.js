console.log('kees client v0.0.1');

class KeesClient {
  constructor(){
    this.nav = this.wrap(document.querySelector('.nav'));
    this.content = this.wrap(document.querySelector('.content'));
    this.storage = {
      currentDevice: undefined,
      devices: [],
      commands: {},
      command_poller: undefined
    };
  }

  // TODO: look into extending standard node with these helpers
  //       to avoid el.el.helper()
  wrap(el){
    return {
      el: el,
      toggle: () => { el.classList.toggle('hide') },
      show:   () => { el.classList.remove('hide') },
      hide:   () => { el.classList.add('hide')    },
      clear:  () => { el.innerHTML = ''           }
    }
  }

  hook(e){
    let hooks = e.el.querySelectorAll('[hook]');
    hooks.forEach(h => {
      let method = h.attributes.hook.value;
      h.addEventListener('click', (el) => {
        this[method].bind(this)(el);
      });
    })
  }

  assemble(root, tmpl){
    let t = document.createElement("template");
    t.innerHTML = tmpl;
    let node = t.content.cloneNode(true);

    root.clear()
    root.el.appendChild(node);
    this.hook(root);
  }

  generateNode(tmpl){
    let t = document.createElement("template");
    t.innerHTML = tmpl;
    let node = t.content.cloneNode(true);
    return node;
  }

  store(key, data){
    this.storage[key] = data;
    localStorage.setItem('auth', JSON.stringify(data));
  }

  retrieve(key){
    let data = localStorage.getItem(key);
    data = JSON.parse(data);
    this.storage[key] = data;
    return data;
  }

  clear(key){
    delete this.storage[key];
    localStorage.removeItem('auth');
  }

  async request(opts){
    let options = {
      method: opts.method || "GET",
    };

    if (this.storage.auth){
      options.headers = { 'X-Kees-JWT-Token': this.storage.auth.jwt.token };
    }

    if (opts.payload){
      options.body = JSON.stringify(opts.payload);
    }

    let resp = await fetch(opts.path, options);

    let data = await resp.json();
    return {http: resp, data: data};
  }


  start(){
    this.storage.auth = this.retrieve('auth');
    this.configureApp();
    this.renderApp();
  }

  async configureApp(){
    if (!this.storage.auth){ return }
  }

  clearApp(){
    this.clear('auth');
    this.start();
  }



  renderApp(){
    if (this.storage.auth){ this.initializeApp() }
    else { this.initializeLogin() }
  }

  initializeApp(){
    this.loggedInNav(this.nav);
    this.refreshDevices();
  }

  initializeLogin(){
    this.loggedOutNav(this.nav);
    this.renderLogin();
  }

  loggedInNav(){
    let tmpl = `
    `;

    this.assemble(nav, tmpl);
  }

  loggedOutNav(nav){
    this.assemble(nav, `
      <div>
        <a href="#login" hook="renderLogin">Login</a> | <a href="#about" hook="renderAbout">About</a>
      </div>
    `);

  }

  loggedInNav(nav){
    this.assemble(nav, `
      <div>
        <a href="#devices" hook="refreshDevices">Devices</a> |
        <a href="#settings" hook="renderSettings">Settings</a> |
        <a href="#about" hook="renderAbout">About</a>
      </div>
    `);
  }

  renderLogin(){
    this.assemble(this.content, `
      <div class="auth row">
        <form class="details">
          <div class="row">
            <label for="name">username</label>
            <input type="text" id="username" autofocus><br>
          </div>

          <div class="row">
            <label for="password">password</label>
            <input type="password" id="password"><br>
          </div>
        </form>

        <div class="row login">
          <button hook="auth">Login</button>
        </div>
        <div class="row loading hide">
          <button><i class="fa-solid fa-compact-disc fa-spin"></i></button>
        </div>
        <div class="errors">
        </div>
      </div><!-- auth -->
    `);
  }

  // TODO: prefix network call function names (or add class?)
  async auth(){
    let auth = this.getAuthInfo();
    // TODO: add event logger

    let login = this.wrap(this.content.el.querySelector(".login"));
    let loading = this.wrap(this.content.el.querySelector(".loading"));
    let errors  = this.wrap(this.content.el.querySelector(".errors"));
    errors.clear();

    login.hide();
    loading.show();

    let resp = await this.request({
      path: '/api/v1/auth',
      method: 'POST',
      payload: auth
    });


    setTimeout(() => {
      this.handleAuth({
        components: {
          login: login,
          loading: loading,
          errors: errors
        },
        resp: resp
      });
    }, 500);
  }

  // TODO: clear login fields if success/password if failure
  handleAuth(data){
    let resp = data.resp;
    if (resp.http.status != 200){
      data.components.loading.hide();
      data.components.login.show();
      let errors = data.components.errors;
      errors.el.innerHTML = resp.data.message;
      errors.show()
    }

    this.store('auth', resp.data.data);
    this.start();
  }

  getAuthInfo(){
    let auth = this.content.el.querySelector(".auth");
    let info = {
      username: auth.querySelector('#username').value,
      password: auth.querySelector('#password').value,
      client: window.navigator.userAgent
    }

    return info;
  }

  clearAuthInfo(){
    let auth = this.content.el.querySelector(".auth");
    auth.querySelector('#username').value = '';
    auth.querySelector('#password').value = '';
  }


  async refreshDevices(){
    // Fetch list of devices and store locally
    let devices = await this.getDevices();
    if (devices == undefined){
      return
    }
    this.storage.devices = devices.data;
    for (const i in this.storage.devices) {
      let device = this.storage.devices[i];
      console.debug(device);
      this.storage.commands[device.id] = [];
      await this.getCommandHistory(device);
    }

    this.storage.currentDevice = this.storage.devices[0];
    this.renderDeviceMain();
  }

  async getDevices(){
    console.log("getting device list");
    // TODO: add event log for device grab
    let resp = await this.request({
      path: '/api/v1/devices'
    });

    // TODO: add generic error handler
    if (resp.http.status != 200){
      if (resp.data.message == 'Invalid JWT'){
        this.clearApp();
        return;
      }
    }

    console.debug({resp});
    return resp;
  }

  async getCommandHistory(device, page = 0, poller = undefined){
    console.info(`getting history for ${device.id}`);
    let endpoint = `/api/v1/devices/${device.id}/commands?page=${page}`;
    let resp = await this.request({
      path: endpoint,
    });

    console.debug(resp);

    this.storage.commands[device.id] = resp.data.data;

    let pending = false;
    for (const i in this.storage.commands[device.id].commands){
      let command = this.storage.commands[device.id].commands[i];
      if (command.status == "pending") {
        pending = true;
      }
    }

    if (pending && (poller || this.storage.command_poller == undefined)){
      this.storage.command_poller = setTimeout( e => {
        this.getCommandHistory(device, page, true)
      }, 500);
    }

    if (!pending) {
      this.storage.command_poller = undefined;
      this.renderDeviceMain();
    }
  }

  renderDeviceMain(){
    let devices = this.storage.devices || [];
    let deviceList = this.renderDeviceList(devices);
    let device = this.renderDeviceInfo(devices[0]);

    this.assemble(this.content, `
      <div class="app row">
        ${deviceList}
        ${device}
      </div>
    `);
  }



  renderCommandHistory(device){
    let base = `
      <div class="row history header">
        <div class="date four columns"><strong>date</strong></div>
        <div class="operation four columns"><strong>Operation</strong></div>
        <div class="status four columns"><strong>Status</strong></div>
      </div>
    `;


    let commands = this.storage.commands[device.id];
    let previous_page = "&nbsp;";
    let next_page = "&nbsp;";

    if (commands.commands.length == commands.page_count){
      next_page = `
        <button hook="commandPagination" data-command-page=${commands.page + 1}>Next »</button>
      `;
    }

    if (commands.page != 0 || (commands.commands.length < commands.page_count)){
      previous_page = `
        <button hook="commandPagination" data-command-page=${commands.page - 1}>« Previous</button>
      `;
    }

    let pagination = `
      <div class="row pagination">
        <div class="six columns prev">${previous_page}</div>
        <div class="six columns prev">${next_page}</div>
      </div>
    `;

    let commandDiv = commands.commands.reduce((all, e) => {
      let pending = (e.status == "pending") ? 'fa-solid fa-circle-notch fa-spin' : '';
      let t = new Date(e.created_at);
      return all + `
        <div class="row history" title="${e.created_at}">
          <div class="created_at four columns">${t.toLocaleDateString()} @ ${t.toLocaleTimeString() }</div>
          <div class="operation four columns">${e.operation}</div>
          <div class="status four columns"><span class="${e.status}">${e.status}</span> <i class="${pending}"></i></div>
        </div>
      `;
    }, base);


    commandDiv += pagination;

    return commandDiv;
  }

  commandPagination(e){
    let page = e.srcElement.dataset.commandPage;
    this.getCommandHistory(this.storage.currentDevice, page = page);
  }

  renderDeviceInfo(device){
    let capabilities = device.capabilities.reduce((all, c) => {
      return all + this.renderCapabilities(device, c);
    }, '');


    let history = this.renderCommandHistory(device);

    let deviceInfo = `
      <div class="device eight columns">
        <div class="row details">
          <div class="row name">${device.name}</div>
          <div class="row last_heartbeat">
            Last Seen: <span class="hb">${ this.fuzzyTime(device.last_heartbeat) }</strong>
          </div>

          <div class="row controller">
            ${device.controller} @ <span class="version">${device.version}</span>
          </div>
        </div>

        <div class="row actions">
          <div class="row capabilities">
            ${capabilities}
          </div>
        </div>

        <div class="row commands">
          <div class="row header">Commands</div>
          <div class="row command_list">
            ${history}
          </div>
        </div>
      </div>

    `;
    return deviceInfo;
  }

  renderCapabilities(device, capability){
    let capabilityBtn = `
      <button
        ${ device.online ? '' : 'disabled' }
        class="action"
        data-device-id=${device.id}
        data-device-action="${capability}"
        hook="deviceAction">
          ${capability}
      </button>
    `;

    return capabilityBtn;
  }

  async deviceAction(e){
    let el = e.srcElement;

    let deviceID = el.dataset.deviceId;
    let action = el.dataset.deviceAction;

    let endpoint = `/api/v1/devices/${deviceID}/commands/${action}`;
    let resp = await this.request({
      path: endpoint,
      method: 'POST'
    });

    await this.getCommandHistory(this.storage.devices[0]);
    this.renderDeviceMain();
  }


  renderDeviceList(devices){
    let deviceElements = devices.map(e => {
      return this.renderDevice(e);
    });


    let deviceList = `
      <div class="devices four columns">
        <div class="row header">Device List</div>
        ${deviceElements}
      </div>
    `;
    return deviceList;
  }

  renderDevice(device){
    let online = device.online;
    return `
      <div class="device row" data-device-id="${device.id}">
        <div class="two columns status ${ online ? 'online' : 'offline'}">
          <i class="fa-solid ${ online ? 'fa-signal' : 'fa-power-off'}"></i>
        </div>
        <div class="eight columns info">
          <div class="row name">${device.name}</div>
          <div class="row controller">${device.controller}</div>
        </div>
        <div class="two columns arrow">
          <i class="fa-solid fa-play"></i>
        </div>
      </div>
    `
  }


  renderAbout(){
    this.assemble(this.content, `
      <div class="about row">
        <p>This is a paragraph about kees</p>
        <p>You can find more at <a href="https://github.com/nhmood/kees">GitHub</a></p>
      </div>
    `);
  }

  // Friendly time
  fuzzyTime(time){
    var delta = (new Date() - (new Date(time))) / 1000;
    delta = Math.round(delta)
    var minute  = 60,
        hour    = minute * 60,
        day     = hour * 24,
        week    = day * 7,
        month   = week * 4,
        year    = month * 12;

    var fuzzy;

    if (delta < 30) {
          fuzzy = 'just now';
    } else if (delta < minute) {
          fuzzy = delta + ' seconds ago';
    } else if (delta < 2 * minute) {
          fuzzy = 'a minute ago'
    } else if (delta < hour) {
          fuzzy = Math.floor(delta / minute) + ' minutes ago';
    } else if (Math.floor(delta / hour) == 1) {
          fuzzy = '1 hour ago.'
    } else if (delta < day) {
          fuzzy = Math.floor(delta / hour) + ' hours ago';
    } else if (delta < day * 2) {
          fuzzy = 'yesterday';
    } else if (delta < week) {
        fuzzy = Math.floor(delta / day) + ' days ago';
    } else if (Math.floor(delta / week) == 1){
      fuzzy = '1 week ago.';
    } else if (delta < month) {
      fuzzy = Math.floor(delta / week) + ' weeks ago';
    } else if (delta < year) {
      fuzzy = Math.floor(delta / month) + ' months ago';
    } else {
      fuzzy = Math.floor(delta / year) + ' years ago';
    }

    return fuzzy;
  }


}

const kees = new KeesClient();

document.addEventListener("DOMContentLoaded", function() {
  kees.start();
});
