import React from 'react';


var Item = React.createClass({
  render() {
    const isActive = this.props.href === (window.location.pathname + '#' + this.props.location.pathname);
    const activeClass = isActive ? 'active' : '';

    return (
      <div className={'nav-item-container nav-item-container-' + this.props.cols}>
        <div className={'nav-item ' + activeClass + ' ' + this.props.className}>
          <a href={this.props.href}>
            <span>{this.props.label}</span>
            <img src={this.props.icon}/>
          </a>
        </div>
      </div>
    );
  }
});

export default Item;
