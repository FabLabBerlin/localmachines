var React = require('react');


var ProductPage = React.createClass({
  render() {
    return (
      <div id="prod">
        <div id="prod-head">
          <h1>EASY LAB</h1>
          <h3>The operating system for your makerspace.</h3>
          <button id="prod-demo-button">Try the demo</button>
        </div>
        <div className="row">
          <div className="col-md-6">
            <h2>Admin & User Webinterface</h2>
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
        </div>
        <div className="row">
          <div className="col-md-6">
            <img src="/machines/assets/img/product/Plug.png"/>
          </div>
          <div className="col-md-6">
            <h2>The Hardware</h2>
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
        </div>
        <div className="row">
          <div className="col-xs-12">
            <h2>Ready to come on board?</h2>
            <p>Ok, you’ll give us a shot? Contact us to become a free Beta Tester.</p>
          </div>
        </div>
        <div className="row">
          <div className="col-md-3"/>
          <div className="col-md-2">
            Send us a mail.
          </div>
          <div className="col-md-2">
            Give us a call.
          </div>
          <div className="col-md-2">
            Drop by.
          </div>
          <div className="col-md-3"/>
        </div>
        <div className="row">
          <div className="col-md-3"/>
          <div className="col-md-2">
            <a href="mailto:info@easylab.io">info@easylab.io</a>
          </div>
          <div className="col-md-2">
            <a href="+4917645839279">+49 176 45839279</a>
          </div>
          <div className="col-md-2">
            <div>Fab Lab Berlin/Makea Industries GmbH</div>
            <div>Prenzlauer Allee 242, 10405 Berlin</div>
          </div>
          <div className="col-md-3"/>
        </div>
      </div>
    );
  }
});

export default ProductPage;
