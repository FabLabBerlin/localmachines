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
            <div className={'col-md-9 ' +
                            'text-xs-center text-md-left'}>
              {text}
            </div>
            <div className="col-md-3"/>
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
            <div className="col-md-3"/>
            <div className={'col-md-9 ' +
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
                 title="Wolf Jeschonnek">
          Wolf is the founder and CEO of Fab Lab Berlin/Makea Industries GmbH.
          After his diploma in Product-Design at the academy of art Berlin
          Weißensee he went out to explore the world of makers and establish
          the first Fab Lab in Berlin. He takes care of the project management
          and product development. 
        </Profile>

        <Profile image="/machines/assets/img/product/phil.jpg"
                 left={false}
                 title="Philip Silva">
          Philip is the software developer in our team. He makes sure that even
          the craziest feature requests become reality. He is constantly aiming
          for a more reliable, convenient and secure system. 
        </Profile>

        <Profile image="/machines/assets/img/product/charlie.jpg"
                 left={true}
                 title="Charlie-Camille Thomas">
          Charlie studied applied art in Paris, and Graphic-Design in Geneva.
          She is now taking care of all the nasty, little details that make up
          a neat and well working interface. Even though EASY LAB is
          whitelabeled, true simplicity is hard work.
        </Profile>

        <Profile image="/machines/assets/img/product/sylwes.jpg"
                 left={false}
                 title="Sylwester Sosnowski">
          Sylwester is Philip's closests counterpart when it comes to technical
          development. He takes care of the hardware development and is
          constantly trying out new ways of interacting with the machines.
        </Profile>

        <Profile image="/machines/assets/img/product/max.jpg"
                 left={true}
                 title="Maximilian Mahal">
          Max is focusing on the work of a Lab Manager and his customers and
          tries to match them in the user experience of EASY LAB. A truly
          iterative process. 
          Your critique is always welcome and valuable input. 
        </Profile>

        <Profile image="/machines/assets/img/product/murat.jpg"
                 left={false}
                 title="Murat Vurucu">
          Murat is a co-founder of Fab Lab Berlin/Makea Industries GmbH and a
          true visionary. He helps other makerspaces to make the most out of
          the features we build for our own lab. He dreams of a worldwide
          EASY LAB network to provide easy access for everyone. 
        </Profile>
      </div>
    );
  }
});

export default Profiles;
