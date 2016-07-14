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
