var _ = require('lodash');
var $ = require('jquery');
var Button = require('../Button');
var Location = require('../../modules/Location');
var React = require('react');
var reactor = require('../../reactor');
var toastr = require('toastr');

// https://github.com/HubSpot/vex/issues/72
var vex = require('vex-js'),
VexDialog = require('vex-js/js/vex.dialog.js');

vex.defaultOptions.className = 'vex-theme-custom';


var Input = React.createClass({
  handleEdit(e) {
    if (this.props.onChange) {
      if (this.isBool()) {
        this.props.onChange(!this.props.value);
      } else {
        this.props.onChange(e.target.value);
      }
    }
  },

  isBool() {
    return this.props.value === false || this.props.value === true;
  },

  render() {
    if (this.isBool()) {
      return (
        <input onChange={this.handleEdit}
               type="checkbox"
               checked={this.props.value}/>
      );
    } else {
      return (
        <input onChange={this.handleEdit}
               value={this.props.value}/>
      );
    }
  }
});


// Plain Table with Create-Read-Update-Delete function
var TableCRUD = React.createClass({

  getInitialState() {
    return {
    };
  },

  handleAdd() {
    this.props.onAdd();
  },

  handleArchive(i) {
    if (!this.state.editRow) {
      return;
    }

    if (this.state.editRow.index !== i) {
      toastr.error('Internal Error.  Please try again later.');
      return;
    }

    var cb = () => {
      const entity = this.props.entities.get(i).toJS();

      const locationId = reactor.evaluateToJS(Location.getters.getLocationId);

      $.ajax({
        url: this.props.updateUrl + '/' + entity.Id + '/archive?location=' + locationId,
        dataType: 'json',
        type: 'PUT',
        contentType: 'application/json; charset=utf-8'
      })
      .done(() => {
        toastr.info('Successfully archived.');

        this.setState({
          editRow: null
        });
        if (this.props.onAfterUpdate) {
          this.props.onAfterUpdate();
        } else {
          toastr.error('Please define onAfterUpdate property');
        }
      })
      .fail(() => {
        toastr.error('Error saving.  Please try again later.');
      });
    };

    VexDialog.buttons.YES.text = 'Yes';
    VexDialog.buttons.NO.text = 'No';
    
    VexDialog.confirm({
      message: 'Do you really want to archive this purchase?',
      callback: confirmed => {
        if (confirmed) {
          cb();
        }
        $('.vex').remove();
        $('body').removeClass('vex-open');
      }
    });
  },

  handleClick(i) {
    if (this.state.editRow) {
      if (this.state.editRow.index !== i) {
        /*eslint-disable no-alert */
        if (this.state.editRow.changes &&
          !window.confirm('Really continue aborting changes?')) {

          return;
        }
        /*eslint-enable no-alert */

        this.setState({
          editRow: null
        });
      }
    } else {
      this.setState({
        editRow: {
          index: i
        }
      });
    }
  },

  handleEdit(i, key, value) {
    var editRow = _.cloneDeep(this.state.editRow);

    if (!editRow.changes) {
      editRow.changes = {};
    }

    editRow.changes[key] = value;

    this.setState({
      editRow: editRow
    });
  },

  handleSave(i) {
    if (!this.state.editRow) {
      return;
    }

    if (this.state.editRow.index !== i) {
      toastr.error('Internal Error.  Please try again later.');
      return;
    }

    const entity = this.props.entities.get(i).toJS();

    _.each(this.state.editRow.changes, (v, k) => {
      entity[k] = v;
    });

    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);

    $.ajax({
      url: this.props.updateUrl + '/' + entity.Id + '?location=' + locationId,
      dataType: 'json',
      type: 'PUT',
      contentType: 'application/json; charset=utf-8',
      data: JSON.stringify(entity)
    })
    .done(() => {
      toastr.info('Successfully saved.');

      this.setState({
        editRow: null
      });
      if (this.props.onAfterUpdate) {
        this.props.onAfterUpdate();
      } else {
        toastr.error('Please define onAfterUpdate property');
      }
    })
    .fail(() => {
      toastr.error('Error saving.  Please try again later.');
    });
  },

  handleSetShowArchived(yes) {
    this.setState({
      showArchived: yes
    });
  },

  render() {
    const showArchived = this.state.showArchived;
    console.log('this.state.showArchive=', showArchived);

    return (
      <div className="container">
        <h1>{this.props.title}</h1>

        <table className="table table-striped table-hover">
          <thead>
            <tr>
              {this.props.fields.map(f => {
                return <th key={f.key}>{f.label}</th>;
              })}
              <th/>
            </tr>
          </thead>

          <tbody>
            {this.props.entities.filter(e => !e.get('Archived') || showArchived)
                                .map((e, i) => {
              const editRow = this.state.editRow &&
                              this.state.editRow.index === i &&
                              e.get('Id') !== 0;

              return (
                <tr key={i} onClick={this.handleClick.bind(this, i)}>
                  {this.props.fields.map(f => {
                    var value = e.get(f.key);

                    if (editRow && f.key !== 'Id') {
                      if (this.state.editRow.changes &&
                          this.state.editRow.changes[f.key]) {

                        value = this.state.editRow.changes[f.key];
                      }

                      return <Input field={f}
                                    value={value}
                                    onChange={this.handleEdit.bind(this, i, f.key)}/>;
                    } else {
                      if (value === true || value === false) {
                        if (value) {
                          return <i className="fa fa-check"/>;
                        } else {
                          return <div/>;
                        }
                      } else {
                        return e.get(f.key);
                      }
                    }
                  }).map((tag, j) => <td key={j}>{tag}</td>)}

                  <td>
                    {editRow ? <i className="fa fa-archive"
                                  onClick={this.handleArchive.bind(this, i)}
                                  style={{cursor: 'pointer', marginRight: '15px'}}/> : null}
                    {editRow ? <i className="fa fa-floppy-o"
                                  onClick={this.handleSave.bind(this, i)}
                                  style={{cursor: 'pointer'}}/> : null}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>

        <div style={{ height: '50px' }}>
          <Button.Annotated id="inv-add-purchase"
                            icon="/machines/assets/img/invoicing/add_purchase.svg"
                            label="Add"
                            onClick={this.handleAdd}/>
        </div>
        <div style={{ height: '50px' }}>
          <Button.Annotated label={showArchived ? 'Hide Archived' : 'Show Archived'}
                            icon="/machines/assets/img/invoicing/add_purchase.svg"
                            onClick={this.handleSetShowArchived.bind(this, !showArchived)}/>
        </div>
      </div>
    );
  }
});

export default TableCRUD;
