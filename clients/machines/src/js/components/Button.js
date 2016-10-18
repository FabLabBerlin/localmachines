var React = require('react');


var Annotated = React.createClass({
  click() {
    this.props.onClick();
  },

  render() {
    const iconStyle = {
      backgroundImage: 'url(' + this.props.icon + ')'
    };

    return (
      <button className="button-annotated"
              id={this.props.id}
              onClick={this.click}>
        <div className="button-annotated-icon"
             style={iconStyle}/>
        <div>{this.props.label}</div>
      </button>
    );
  }
});

export default {
  Annotated
};
