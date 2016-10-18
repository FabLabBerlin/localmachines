var React = require('react');


var Annotated = React.createClass({
  click(e) {
    this.props.onClick(e);
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
        <div className="button-annotated-label">
          {this.props.label}
        </div>
      </button>
    );
  }
});

export default {
  Annotated
};
