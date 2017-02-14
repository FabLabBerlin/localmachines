var $ = require('jquery');
import React from 'react';


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


var Tiny = React.createClass({
  componentDidMount() {
    $(this.i).tooltip();
  },

  handleClick(e) {
    if (this.props.onClick) {
      this.props.onClick(e);
    }
  },

  render() {
    var cls = this.props.faClassName + ' ';
    cls += this.props.className;

    return (
      <i className={cls}
         aria-hidden="true"
         data-toggle="tooltip"
         data-placement="top"
         data-toggle="tooltip"
         onClick={this.handleClick}
         ref={i => { this.i = i; }}
         style={this.props.style}
         title={this.props.title}/>
    );
  }
});

export default {
  Annotated,
  Tiny
};
