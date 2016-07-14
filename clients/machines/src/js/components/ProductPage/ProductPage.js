var $ = require('jquery');
var GlobalActions = require('../../actions/GlobalActions');
var Profiles = require('./Profiles');
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


// Footer Call-To-Action
var FooterCTA = React.createClass({
  render() {
    return (
      <div className="prod-footer-cta-container">
        <div id={this.props.id}
             className="prod-footer-cta">
          <img src={this.props.image}/>
                {this.props.text}
        </div>
        <div>
          {this.props.children}
        </div>
      </div>
    );
  }
});


var Footer = React.createClass({
  render() {
    return (
      <div id="prod-footer" className="row">
        <div className="col-md-6 text-md-left text-xs-center">
          Easy Lab is a product of Makea Industries GmbH. © 2016
        </div>
        <div className="col-md-6 text-md-right text-xs-center">
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

        <section id="prod-about">
          <div className="row">
            <h3 className="prod-section-title">Two sides, one system.</h3>
            <div className="col-md-6 col-md-push-6">
              <img className="prod-section-image"
                   src="/machines/assets/img/product/PhoneLaptop.png"/>
            </div>
            <div className="col-md-6 col-md-pull-6">
              <h2 className="prod-section-subtitle">A clean web interface.</h2>
              <p>
                EASY LAB is built as a web application on purpose, providing
                trouble free usage no matter what kind of operating system you
                or your customers are working with. It is responsive, of course.
                Depending on the type of account you use to login, you will
                either see the standard user interface or the admin version with
                all the managing options you need to run a makerspace.
              </p>
              <p>
                The main feature for your customers is to activate the machines.
                EASY LAB is keeping track of these activations, making it
                convenient for you to bill them, and providing a clear overview
                of your customer’s spendings. No bad surprises. 
              </p>
              <p>
                If your members are facing serious deadlines, the reservation
                rules come in handy. Allowing them to book machines for specific
                timeframes in order to have privileged access. You can charge for
                that, in order to keep a good culture. Reservations are visible
                for all other users to plan accordingly, of course. 
              </p>
            </div>
          </div>

          <div className="row">
            <div className="col-md-6">
              <img className="prod-section-image"
                   src="/machines/assets/img/product/Screens.jpg"/>
            </div>
            <div className="col-md-6">
              <p>
                We think face-to-face communication still provides the most
                bandwidth and is the fastest way to solve things. Even more
                when you want to cater for a lively community. However,
                EASY LAB provides useful, but therefore not overly
                exaggerating, communicational features, too. Users can easily
                report a machine failure and provide general feedback about
                billing, technical issues and so on. 
              </p>
              <p>
                In the admin version you can of course add and delete machines,
                change their specs and pricing, but also set them into
                maintenance mode if they need a little rest. 
                You can keep track of your user’s permissions to run the
                machines, which will let you sleep assured that only trained
                users are operating your valuable infrastructure. 
                Of course, you can set all the rules for reservations,
                membership models here, too. 
              </p>
              <p>
                Our latest trick is to provide you with a neat automatic
                invoicing feature, which will save you a lot of time at the end
                of the month. 
              </p>
            </div>
          </div>

          <div className="row">
            <div className="col-md-6 col-md-push-6">
              <img className="prod-section-image"
                   src="/machines/assets/img/product/Plug.png"/>
            </div>
            <div className="col-md-6 col-md-pull-6">
              <h2 className="prod-section-subtitle">A reliable hardware setup.</h2>
              <p>
                EASY LAB relies on wifi enabled powerswitches which turn your
                machines on and off. At the heart of it all is the gateway
                which comes as an easy to setup Image for a Raspberry Pi 3.
                Connecting a new machine to your setup is a one time process of
                round about 10 minutes. Placing the interaction outside of the
                machine itself means more flexibility and mobility for your
                space, making it easy to rearrange your infrastructure as you
                are growing or moving in a new location. It even supports
                mobile setups. Built on a Raspberry Pi, allows us to further
                develop more advanced machine interactions in the near future.
                Think: maintenance counters, job specific machine interactions,
                remote monitoring etc. 
              </p>
            </div>
          </div>
        </section>

        <div className="prod-section-separator">
          If you want to explore all of the features of EASY LAB, why not
          become <span className="prod-section-separator-link" onClick={this.clickContact}>a free Beta-tester</span>?
        </div>

        <Profiles/>

        <div id="prod-contact">
          <div id="prod-ready-to-board" className="row">
            <div className="col-xs-12">
              <h2 className="prod-section-title text-center">
                EASY LAB is work in progress.
                <br/>
                Try the free beta.
              </h2>
              <p className="text-center">
                Get in contact if you want to  use EASY LAB in your Makerspace. 
              </p>
              <p className="text-center">
                Did i mention the price point? - Zero.
              </p>
            </div>
          </div>

          <div className="row">
            <div className="col-md-1"/>
            <div className="col-md-3 text-center">
              <FooterCTA id="prod-footer-cta-send-mail"
                         image="/machines/assets/img/product/send_mail.svg"
                         text="Send us a mail.">
                <a href="mailto:easylab@makea.org">easylab@makea.org</a>
              </FooterCTA>
            </div>
            <div className="col-md-4 text-center">
              <FooterCTA id="prod-footer-cta-call"
                         image="/machines/assets/img/product/call.svg"
                         text="Give us a call.">
                <a href="tel:+4917645839279">+49 176 45839279</a>
              </FooterCTA>
            </div>
            <div className="col-md-3 text-center">
              <FooterCTA id="prod-footer-cta-drop-by"
                         image="/machines/assets/img/product/drop_by.svg"
                         text="Drop by.">
                <a href="https://goo.gl/maps/k1ksF5AjDUD2">
                  <div>Fab Lab Berlin/Makea Industries GmbH</div>
                  <div>Prenzlauer Allee 242, 10405 Berlin</div>
                </a>
              </FooterCTA>
            </div>
            <div className="col-md-1"/>
          </div>
        </div>
        <Footer/>
      </div>
    );
  }
});

export default ProductPage;
