console.log('kees client v0.0.1');

class KeesClient {
  constructor(){
    this.session = undefined;
    this.nav = this.wrap(document.querySelector('.nav'));
    this.content = this.wrap(document.querySelector('.content'));
  }

  wrap(el){
    return {
      el: el,
      toggle: () => { el.classList.toggle('hide') },
      show:   () => { el.classList.remove('hide') },
      hide:   () => { el.classlist.add('hide')    },
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


  start(){
    this.renderNav();
    this.renderContent();
  }

  renderNav(){
    if (this.session){
    } else {
      this.loggedOutNav(this.nav)
    }
  }

  renderContent(){
    if (this.session){
    } else {
      this.renderLogin();
    }
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

  renderLogin(){
    this.assemble(this.content, `
      <div class="auth row">
        <form class="details">
          <div class="row">
            <label for="name">username</label>
            <input type="text" id="username"><br>
          </div>

          <div class="row">
            <label for="password">password</label>
            <input type="password" id="password"><br>
          </div>

        </form><!-- auth -->
        <button hook="auth">Login</button>
      </div>
    `);
  }


  async request(opts){
    let resp = await fetch(opts.path, {
      method: opts.method || "GET",
      body: JSON.stringify(opts.payload)
    });

    let data = await resp.json();
    return data;
  }

  async auth(){
    let auth = this.getAuthInfo();
    // TODO: add event logger

    let resp = await this.request({
      path: '/api/v1/auth',
      method: 'POST',
      payload: auth
    });

    console.log(resp);
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

  renderAbout(){
    this.assemble(this.content, `
      <div class="about row hide">
        <p>This is a paragraph about kees</p>
        <p>You can find more at <a href="https://github.com/nhmood/kees">GitHub</a></p>
      </div>
    `);
  }


}

const kees = new KeesClient();

document.addEventListener("DOMContentLoaded", function() {
  kees.start();
});
