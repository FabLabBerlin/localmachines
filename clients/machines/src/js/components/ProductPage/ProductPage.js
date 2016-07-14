var $ = require('jquery');
var GlobalActions = require('../../actions/GlobalActions');
var React = require('react');
var toastr = require('../../toastr');


var Subscribe = React.createClass({
  render() {
    return (
      <div id="prod-subscribe">
        <span id="prod-subscribe-at">
          @
        </span>
        <span id="prod-subscribe-email-container">
          <input id="prod-subscribe-email"
                 autoCorrect="off"
                 autoCapitalize="off"
                 autoFocus
                 ref="email"
                 type="text"/>
        </span>
        <button id="prod-subscribe-button"
                onClick={this.subscribe}
                type="button">
          Subscribe
        </button>
      </div>
    );
  },

  subscribe() {
    const email = this.refs.email.getDOMNode().value;

    if (email && email.length > 3) {
      GlobalActions.performSubscribeNewsletter(email);
    } else {
      toastr.error('Please provide a valid E-Mail address.');
    }
  }
});


var Profile = React.createClass({
  render() {
    const left = this.props.left;
    const direction = left ? 'left' : 'right';
    const img = <img src={this.props.image}/>;

    const text = (
      <div className="prod-profile">
        <div className="row">
          <h3 className={'prod-profile-title ' + 
                         (left ? '' : 'pull-right ') +
                         'prod-profile-text-' +
                         'prod-profile-title-' + direction}>
            {this.props.title}
          </h3>
        </div>
        <div className={'row prod-profile-text-' + direction}>
          <p>
            {this.props.children}
          </p>
        </div>
      </div>
    );

    return (
      <div className="row">
        <div className={left ? 'col-md-3' : 'col-md-9'}>
          {left ? img : text}
        </div>
        <div className={left ? 'col-md-9' : 'col-md-3'}>
          {left ? text : img}
        </div>
      </div>
    );
  }
});


// Footer Call-To-Action
var FooterCTA = React.createClass({
  render() {
    return (
      <div id={this.props.id}
           className="prod-footer-cta">
        <img src={this.props.image}/>
              {this.props.text}
      </div>
    );
  }
});


var Footer = React.createClass({
  render() {
    return (
      <div id="prod-footer" className="row">
        <div className="col-md-6">
          Easy Lab is a product of Makea Industries GmbH. © 2016
        </div>
        <div className="col-md-6 text-right">
          <a href="https://fablab.berlin/de/content/2-Impressum">
            Imprint
          </a>
        </div>
      </div>
    );
  }
});


var ProductPage = React.createClass({
  click(id) {
    $('html, body').animate({
      scrollTop: $(id).offset().top
    }, 500);
  },

  clickAbout() {
    this.click('#prod-about');
  },

  clickContact() {
    this.click('#prod-contact');
  },

  clickLogin() {
    window.location.href = 'https://easylab.io';
  },

  clickTeam() {
    this.click('#prod-team');
  },

  render() {
    return (
      <div id="prod" className="container-fluid">
        <div id="prod-nav row">
          <div className="col-md-2">
            <img id="prod-makea-logo"
                 src="/machines/assets/img/product/Makea_Logo.png"/>
          </div>
          <div className="col-md-10 hidden-xs hidden-sm text-right">
            <button className="prod-nav-button"
                    onClick={this.clickAbout}
                    type="button">About</button>
            <button className="prod-nav-button"
                    onClick={this.clickTeam}
                    type="button">Team</button>
            <button className="prod-nav-button"
                    onClick={this.clickContact}
                    type="button">Contact</button>
            <button id="prod-nav-login"
                    className="prod-nav-button"
                    onClick={this.clickLogin}
                    type="button">Login</button>
          </div>
        </div>
        <div id="prod-head">
          <h1 id="prod-title">EASY LAB</h1>
          <h3 id="prod-subtitle">
            <p>The operating system for your makerspace.</p>
            <p>Make your lab work instead of micromanaging it.</p>
            <p id="prod-subscribe-title">
              Subscribe to our mailing list to receive the latest info:
            </p>
          </h3>
          <Subscribe/>
        </div>

        <section id="prod-about" className="row">
          <div className="col-md-6">
            <h2 className="prod-section-title">Admin & User Webinterface</h2>
            <p>
              Easy Lab consists mainly of a webinterface which provides
              blablabla…Lorem ipsum dolor sit amet, consetetur sadipscing
              elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore
              magna aliquyam erat, sed diam voluptua. At vero eos et accusam et
              justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea
              takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor
              sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod
              tempor invidunt ut labore et dolore magna aliquyam erat, sed diam
              voluptua. At vero eos et accusam et justo duo dolores et ea rebum.
              Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum
              dolor sit amet.
            </p>
          </div>
          <div className="col-md-6">
            <img src="/machines/assets/img/product/PhoneLaptop.png"/>
          </div>
        </section>

        <section className="row">
          <div className="col-md-6">
            <img src="/machines/assets/img/product/Plug.png"/>
          </div>
          <div className="col-md-6">
            <h2 className="prod-section-title">The Hardware</h2>
            <p>
              You connect your machines via WiFi enabled power switches and a
              Raspberry Pi 3…Lorem ipsum dolor sit amet, consetetur sadipscing
              elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore
              magna aliquyam erat, sed diam voluptua. At vero eos et accusam et
              justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea
              takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum
              dolor sit amet, consetetur sadipscing elitr, sed diam nonumy
              eirmod tempor invidunt ut labore et dolore magna aliquyam erat,
              sed diam voluptua. At vero eos et accusam et justo duo dolores et
              ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est
              Lorem ipsum dolor sit amet.
            </p>
          </div>
        </section>

        <div id="prod-team" className="row">
          <div className="col-xs-12">
            <h2 className="prod-section-title text-center">
              <div>We are running our own Fab Lab</div>
              <div>Oh boy, we know what hustling means.</div>
            </h2>
          </div>
        </div>

        <Profile image="/machines/assets/img/product/Wolf.jpg"
                 left={true}
                 title="The captain">
          Wolf Jeschonnek, founder and CEO of Fab Lab Berlin/Makea
          Industries Gmbh, is guiding this ship through […] He’s
          a real sailer, btw. No kidding.
        </Profile>

        <Profile image="/machines/assets/img/product/phil.jpg"
                 left={false}
                 title="The cook">
          A crew is only as good as their meals. Thanks to Philip
          Silva, who’s mixing up the finest compositions of code
          in the seven seas, we’re good to go that extra mile.
          Philip is the main developer in our team and makes sure
          that even the craziest feature requests become reality.
        </Profile>

        <Profile image="/machines/assets/img/product/charlie.jpg"
                 left={true}
                 title="OC Design">
          Charlie-Camille Thomas is our officer commanding everything
          about the look and feel of Easy Lab. You think, well that’s
          easy because Easy Lab is fully whitelabel. Sorry Bro, but
          true simplicity is really hard work.
        </Profile>

        <Profile image="/machines/assets/img/product/sylwes.jpg"
                 left={false}
                 title="Hardware Guru">
          Sylwester Sosnowski operates our engines. He makes sure that
          the Easy Lab hardware is running as precise and inconspicuously
          like a german u-boot.
        </Profile>

        <Profile image="/machines/assets/img/product/max.jpg"
                 left={true}
                 title="The helmsman">
          No cruise without a profound helmsman, who knows how to steer
          a ship. Maximilian Mahal is rethinking every single grip in the
          workflow of a hard working Lab Manager. He’s always up to chat
          with you…if you can stand him spinning a yarn…
        </Profile>


        <div id="prod-ready-to-board" className="row">
          <div className="col-xs-12">
            <h2 className="prod-section-title text-center">
              Ready to come on board?
            </h2>
            <p className="text-center">
              Ok, you’ll give us a shot? Contact us to become a free Beta Tester.
            </p>
          </div>
        </div>
        <div id="prod-contact" className="row">
          <div className="col-md-1"/>
          <div className="col-md-3 text-center">
            <FooterCTA id="prod-footer-cta-send-mail"
                       image="/machines/assets/img/product/send_mail.svg"
                       text="Send us a mail."/>
          </div>
          <div className="col-md-4 text-center">
            <FooterCTA id="prod-footer-cta-call"
                       image="/machines/assets/img/product/call.svg"
                       text="Give us a call."/>
          </div>
          <div className="col-md-3 text-center">
            <FooterCTA id="prod-footer-cta-drop-by"
                       image="/machines/assets/img/product/drop_by.svg"
                       text="Drop by."/>
          </div>
          <div className="col-md-1"/>
        </div>
        <div className="row">
          <div className="col-md-1"/>
          <div className="col-md-3 text-center">
            <a href="mailto:info@easylab.io">info@easylab.io</a>
          </div>
          <div className="col-md-4 text-center">
            <a href="tel:+4917645839279">+49 176 45839279</a>
          </div>
          <div className="col-md-3 text-center">
            <a href="https://goo.gl/maps/k1ksF5AjDUD2">
              <div>Fab Lab Berlin/Makea Industries GmbH</div>
              <div>Prenzlauer Allee 242, 10405 Berlin</div>
            </a>
          </div>
          <div className="col-md-1"/>
        </div>
        <Footer/>
      </div>
    );
  }
});

export default ProductPage;
