var React = require('react');


var Profile = React.createClass({
  render() {
    const left = this.props.left;
    const direction = left ? 'left' : 'right';
    const img = <img src={this.props.image}/>;

    const text = (
      <div className="prod-profile">
        <div className="row">
          <h3 className={'prod-profile-title ' + 
                         ' prod-profile-title-' + direction +
                         ' text-xs-center ' +
                         (left ? 'text-md-left' : 'text-md-right')}>
            {this.props.title}
          </h3>
        </div>
        <div className={'row prod-profile-text ' +
                        ' prod-profile-text-' + direction +
                        ' text-xs-center' +
                        ' text-md-left'}>
          <p>
            {this.props.children}
          </p>
        </div>
      </div>
    );

    if (left) {
      return (
        <div className="row">
          <div className={'col-md-3 ' +
                          'text-xs-center text-md-left'}>
            {img}
          </div>
          <div className="col-md-9">
            <div className={'' +
                            'text-xs-center text-md-left'}>
              {text}
            </div>
          </div>
        </div>
      );
    } else {
      return (
        <div className="row">
          <div className={'col-md-3 col-md-push-9 ' +
                          'text-xs-center text-md-left'}>
            {img}
          </div>
          <div className="col-md-9 col-md-pull-3">
            <div className={'' +
                            'text-xs-center text-md-left'}>
              {text}
            </div>
          </div>
        </div>
      );
    }
  }
});


var Profiles = React.createClass({
  render() {
    return (
      <div>
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
      </div>
    );
  }
});

export default Profiles;
