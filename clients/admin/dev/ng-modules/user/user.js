(function(){

'use strict';
var app = angular.module('fabsmith.admin.user', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/users/:userId', {
    templateUrl: 'ng-modules/user/user.html',
    controller: 'UserCtrl'
  });
}]); // app.config

app.controller('UserCtrl',
 ['$scope', '$routeParams', '$http', '$location', '$cookies', 'randomToken', 'api',
 function($scope, $routeParams, $http, $location, $cookies, randomToken, api) {

  $scope.locationId = parseInt($cookies.get('location'));
  $scope.globalConfigVisible = $scope.locationId === 1;
  $scope.userMembershipsHistoryVisible = false;

  var countryCodesJSON = '[{"Name":"Afghanistan","Code":"AF"},{"Name":"Åland Islands","Code":"AX"},{"Name":"Albania","Code":"AL"},{"Name":"Algeria","Code":"DZ"},{"Name":"American Samoa","Code":"AS"},{"Name":"Andorra","Code":"AD"},{"Name":"Angola","Code":"AO"},{"Name":"Anguilla","Code":"AI"},{"Name":"Antarctica","Code":"AQ"},{"Name":"Antigua and Barbuda","Code":"AG"},{"Name":"Argentina","Code":"AR"},{"Name":"Armenia","Code":"AM"},{"Name":"Aruba","Code":"AW"},{"Name":"Australia","Code":"AU"},{"Name":"Austria","Code":"AT"},{"Name":"Azerbaijan","Code":"AZ"},{"Name":"Bahamas","Code":"BS"},{"Name":"Bahrain","Code":"BH"},{"Name":"Bangladesh","Code":"BD"},{"Name":"Barbados","Code":"BB"},{"Name":"Belarus","Code":"BY"},{"Name":"Belgium","Code":"BE"},{"Name":"Belize","Code":"BZ"},{"Name":"Benin","Code":"BJ"},{"Name":"Bermuda","Code":"BM"},{"Name":"Bhutan","Code":"BT"},{"Name":"Bolivia, Plurinational State of","Code":"BO"},{"Name":"Bonaire, Sint Eustatius and Saba","Code":"BQ"},{"Name":"Bosnia and Herzegovina","Code":"BA"},{"Name":"Botswana","Code":"BW"},{"Name":"Bouvet Island","Code":"BV"},{"Name":"Brazil","Code":"BR"},{"Name":"British Indian Ocean Territory","Code":"IO"},{"Name":"Brunei Darussalam","Code":"BN"},{"Name":"Bulgaria","Code":"BG"},{"Name":"Burkina Faso","Code":"BF"},{"Name":"Burundi","Code":"BI"},{"Name":"Cambodia","Code":"KH"},{"Name":"Cameroon","Code":"CM"},{"Name":"Canada","Code":"CA"},{"Name":"Cape Verde","Code":"CV"},{"Name":"Cayman Islands","Code":"KY"},{"Name":"Central African Republic","Code":"CF"},{"Name":"Chad","Code":"TD"},{"Name":"Chile","Code":"CL"},{"Name":"China","Code":"CN"},{"Name":"Christmas Island","Code":"CX"},{"Name":"Cocos (Keeling) Islands","Code":"CC"},{"Name":"Colombia","Code":"CO"},{"Name":"Comoros","Code":"KM"},{"Name":"Congo","Code":"CG"},{"Name":"Congo, the Democratic Republic of the","Code":"CD"},{"Name":"Cook Islands","Code":"CK"},{"Name":"Costa Rica","Code":"CR"},{"Name":"Côte d\'Ivoire","Code":"CI"},{"Name":"Croatia","Code":"HR"},{"Name":"Cuba","Code":"CU"},{"Name":"Curaçao","Code":"CW"},{"Name":"Cyprus","Code":"CY"},{"Name":"Czech Republic","Code":"CZ"},{"Name":"Denmark","Code":"DK"},{"Name":"Djibouti","Code":"DJ"},{"Name":"Dominica","Code":"DM"},{"Name":"Dominican Republic","Code":"DO"},{"Name":"Ecuador","Code":"EC"},{"Name":"Egypt","Code":"EG"},{"Name":"El Salvador","Code":"SV"},{"Name":"Equatorial Guinea","Code":"GQ"},{"Name":"Eritrea","Code":"ER"},{"Name":"Estonia","Code":"EE"},{"Name":"Ethiopia","Code":"ET"},{"Name":"Falkland Islands (Malvinas)","Code":"FK"},{"Name":"Faroe Islands","Code":"FO"},{"Name":"Fiji","Code":"FJ"},{"Name":"Finland","Code":"FI"},{"Name":"France","Code":"FR"},{"Name":"French Guiana","Code":"GF"},{"Name":"French Polynesia","Code":"PF"},{"Name":"French Southern Territories","Code":"TF"},{"Name":"Gabon","Code":"GA"},{"Name":"Gambia","Code":"GM"},{"Name":"Georgia","Code":"GE"},{"Name":"Germany","Code":"DE"},{"Name":"Ghana","Code":"GH"},{"Name":"Gibraltar","Code":"GI"},{"Name":"Greece","Code":"GR"},{"Name":"Greenland","Code":"GL"},{"Name":"Grenada","Code":"GD"},{"Name":"Guadeloupe","Code":"GP"},{"Name":"Guam","Code":"GU"},{"Name":"Guatemala","Code":"GT"},{"Name":"Guernsey","Code":"GG"},{"Name":"Guinea","Code":"GN"},{"Name":"Guinea-Bissau","Code":"GW"},{"Name":"Guyana","Code":"GY"},{"Name":"Haiti","Code":"HT"},{"Name":"Heard Island and McDonald Islands","Code":"HM"},{"Name":"Holy See (Vatican City State)","Code":"VA"},{"Name":"Honduras","Code":"HN"},{"Name":"Hong Kong","Code":"HK"},{"Name":"Hungary","Code":"HU"},{"Name":"Iceland","Code":"IS"},{"Name":"India","Code":"IN"},{"Name":"Indonesia","Code":"ID"},{"Name":"Iran, Islamic Republic of","Code":"IR"},{"Name":"Iraq","Code":"IQ"},{"Name":"Ireland","Code":"IE"},{"Name":"Isle of Man","Code":"IM"},{"Name":"Israel","Code":"IL"},{"Name":"Italy","Code":"IT"},{"Name":"Jamaica","Code":"JM"},{"Name":"Japan","Code":"JP"},{"Name":"Jersey","Code":"JE"},{"Name":"Jordan","Code":"JO"},{"Name":"Kazakhstan","Code":"KZ"},{"Name":"Kenya","Code":"KE"},{"Name":"Kiribati","Code":"KI"},{"Name":"Korea, Democratic People\'s Republic of","Code":"KP"},{"Name":"Korea, Republic of","Code":"KR"},{"Name":"Kuwait","Code":"KW"},{"Name":"Kyrgyzstan","Code":"KG"},{"Name":"Lao People\'s Democratic Republic","Code":"LA"},{"Name":"Latvia","Code":"LV"},{"Name":"Lebanon","Code":"LB"},{"Name":"Lesotho","Code":"LS"},{"Name":"Liberia","Code":"LR"},{"Name":"Libya","Code":"LY"},{"Name":"Liechtenstein","Code":"LI"},{"Name":"Lithuania","Code":"LT"},{"Name":"Luxembourg","Code":"LU"},{"Name":"Macao","Code":"MO"},{"Name":"Macedonia, the Former Yugoslav Republic of","Code":"MK"},{"Name":"Madagascar","Code":"MG"},{"Name":"Malawi","Code":"MW"},{"Name":"Malaysia","Code":"MY"},{"Name":"Maldives","Code":"MV"},{"Name":"Mali","Code":"ML"},{"Name":"Malta","Code":"MT"},{"Name":"Marshall Islands","Code":"MH"},{"Name":"Martinique","Code":"MQ"},{"Name":"Mauritania","Code":"MR"},{"Name":"Mauritius","Code":"MU"},{"Name":"Mayotte","Code":"YT"},{"Name":"Mexico","Code":"MX"},{"Name":"Micronesia, Federated States of","Code":"FM"},{"Name":"Moldova, Republic of","Code":"MD"},{"Name":"Monaco","Code":"MC"},{"Name":"Mongolia","Code":"MN"},{"Name":"Montenegro","Code":"ME"},{"Name":"Montserrat","Code":"MS"},{"Name":"Morocco","Code":"MA"},{"Name":"Mozambique","Code":"MZ"},{"Name":"Myanmar","Code":"MM"},{"Name":"Namibia","Code":"NA"},{"Name":"Nauru","Code":"NR"},{"Name":"Nepal","Code":"NP"},{"Name":"Netherlands","Code":"NL"},{"Name":"New Caledonia","Code":"NC"},{"Name":"New Zealand","Code":"NZ"},{"Name":"Nicaragua","Code":"NI"},{"Name":"Niger","Code":"NE"},{"Name":"Nigeria","Code":"NG"},{"Name":"Niue","Code":"NU"},{"Name":"Norfolk Island","Code":"NF"},{"Name":"Northern Mariana Islands","Code":"MP"},{"Name":"Norway","Code":"NO"},{"Name":"Oman","Code":"OM"},{"Name":"Pakistan","Code":"PK"},{"Name":"Palau","Code":"PW"},{"Name":"Palestine, State of","Code":"PS"},{"Name":"Panama","Code":"PA"},{"Name":"Papua New Guinea","Code":"PG"},{"Name":"Paraguay","Code":"PY"},{"Name":"Peru","Code":"PE"},{"Name":"Philippines","Code":"PH"},{"Name":"Pitcairn","Code":"PN"},{"Name":"Poland","Code":"PL"},{"Name":"Portugal","Code":"PT"},{"Name":"Puerto Rico","Code":"PR"},{"Name":"Qatar","Code":"QA"},{"Name":"Réunion","Code":"RE"},{"Name":"Romania","Code":"RO"},{"Name":"Russian Federation","Code":"RU"},{"Name":"Rwanda","Code":"RW"},{"Name":"Saint Barthélemy","Code":"BL"},{"Name":"Saint Helena, Ascension and Tristan da Cunha","Code":"SH"},{"Name":"Saint Kitts and Nevis","Code":"KN"},{"Name":"Saint Lucia","Code":"LC"},{"Name":"Saint Martin (French part)","Code":"MF"},{"Name":"Saint Pierre and Miquelon","Code":"PM"},{"Name":"Saint Vincent and the Grenadines","Code":"VC"},{"Name":"Samoa","Code":"WS"},{"Name":"San Marino","Code":"SM"},{"Name":"Sao Tome and Principe","Code":"ST"},{"Name":"Saudi Arabia","Code":"SA"},{"Name":"Senegal","Code":"SN"},{"Name":"Serbia","Code":"RS"},{"Name":"Seychelles","Code":"SC"},{"Name":"Sierra Leone","Code":"SL"},{"Name":"Singapore","Code":"SG"},{"Name":"Sint Maarten (Dutch part)","Code":"SX"},{"Name":"Slovakia","Code":"SK"},{"Name":"Slovenia","Code":"SI"},{"Name":"Solomon Islands","Code":"SB"},{"Name":"Somalia","Code":"SO"},{"Name":"South Africa","Code":"ZA"},{"Name":"South Georgia and the South Sandwich Islands","Code":"GS"},{"Name":"South Sudan","Code":"SS"},{"Name":"Spain","Code":"ES"},{"Name":"Sri Lanka","Code":"LK"},{"Name":"Sudan","Code":"SD"},{"Name":"Suriname","Code":"SR"},{"Name":"Svalbard and Jan Mayen","Code":"SJ"},{"Name":"Swaziland","Code":"SZ"},{"Name":"Sweden","Code":"SE"},{"Name":"Switzerland","Code":"CH"},{"Name":"Syrian Arab Republic","Code":"SY"},{"Name":"Taiwan, Province of China","Code":"TW"},{"Name":"Tajikistan","Code":"TJ"},{"Name":"Tanzania, United Republic of","Code":"TZ"},{"Name":"Thailand","Code":"TH"},{"Name":"Timor-Leste","Code":"TL"},{"Name":"Togo","Code":"TG"},{"Name":"Tokelau","Code":"TK"},{"Name":"Tonga","Code":"TO"},{"Name":"Trinidad and Tobago","Code":"TT"},{"Name":"Tunisia","Code":"TN"},{"Name":"Turkey","Code":"TR"},{"Name":"Turkmenistan","Code":"TM"},{"Name":"Turks and Caicos Islands","Code":"TC"},{"Name":"Tuvalu","Code":"TV"},{"Name":"Uganda","Code":"UG"},{"Name":"Ukraine","Code":"UA"},{"Name":"United Arab Emirates","Code":"AE"},{"Name":"United Kingdom","Code":"GB"},{"Name":"United States","Code":"US"},{"Name":"United States Minor Outlying Islands","Code":"UM"},{"Name":"Uruguay","Code":"UY"},{"Name":"Uzbekistan","Code":"UZ"},{"Name":"Vanuatu","Code":"VU"},{"Name":"Venezuela, Bolivarian Republic of","Code":"VE"},{"Name":"Viet Nam","Code":"VN"},{"Name":"Virgin Islands, British","Code":"VG"},{"Name":"Virgin Islands, U.S.","Code":"VI"},{"Name":"Wallis and Futuna","Code":"WF"},{"Name":"Western Sahara","Code":"EH"},{"Name":"Yemen","Code":"YE"},{"Name":"Zambia","Code":"ZM"},{"Name":"Zimbabwe","Code":"ZW"}]';
  $scope.countryCodes = JSON.parse(countryCodesJSON);

  // Init scope variables
  var pickadateOptions = {
    format: 'yyyy-mm-dd'
  };
  $('.datepicker').pickadate(pickadateOptions);
  $scope.locations = [];
  $scope.userLocations = [];
  $scope.user = { Id: $routeParams.userId };
  $scope.userMachines = [];
  $scope.userMemberships = [];

  $scope.loadFastbillTemplates = function() {
    $http({
      method: 'GET',
      url: '/api/settings/fastbill_templates',
      params: {
        ac: new Date().getTime(),
        location: $cookies.get('location')
      }
    })
    .success(function(fastbillTemplates) {
      $scope.fastbillTemplates = fastbillTemplates;
    })
    .error(function(data, status) {
      toastr.error('Failed to load fastbill templates');
    });
  };

  $scope.loadLocations = function() {
    $http({
      method: 'GET',
      url: '/api/locations'
    })
    .success(function(locations) {
      $scope.locations = locations;
      _.each(locations, function(l) {
        if (l.Id === $scope.locationId) {
          $scope.location = l;

          return false;
        }
      });
  console.log('$scope.location=', $scope.location);
    })
    .error(function(data, status) {
      toastr.error('Failed to load locations data');
    });
  };

  $scope.loadUserLocations = function() {
    $http({
      method: 'GET',
      url: '/api/users/' + $scope.user.Id + '/locations'
    })
    .success(function(userLocations) {
      $scope.userLocations = userLocations;
      _.each($scope.userLocations, function(ul) {
        if (ul.LocationId === $scope.locationId) {
          $scope.userLocation = _.extend({}, ul);
        }
      });
    })
    .error(function(data, status) {
      toastr.error('Failed to load user locations data');
    });
  };

  $scope.loadUserData = function() {
    $http({
      method: 'GET',
      url: '/api/users/' + $scope.user.Id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(user) {
      if (user.ClientId <= 0) {
        user.ClientId = '';
      }
      user.FastbillTemplateId = String(user.FastbillTemplateId);
      $scope.user = user;
      $scope.loadAvailableMachines();
      $scope.getUserMemberships(); // This part can happen assync
    })
    .error(function(data, status) {
      toastr.error('Failed to load user data');
    });
  };

  // It should go like this:
  //
  // 1. Load user data
  // 2. Load available machines
  // 3. Load user machine permissions
  // 4. Load available memberships
  // 5. Load user memberships

  $scope.loadFastbillTemplates();
  $scope.loadLocations();
  $scope.loadUserLocations();
  $scope.loadUserData();

  $scope.loadAvailableMachines = function() {
    
    api.loadCategories(function(categories) {
      $scope.categories = categories;
      $scope.loadUserCategoryPermissions($scope.getAvailableMemberships);
    });
  };

  $scope.loadUserCategoryPermissions = function(callback) {
    $http({
      method: 'GET',
      url: '/api/users/' + $scope.user.Id + '/permissions',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(userCategories) {

      $scope.userCategories = userCategories;

      if (callback) {
        callback();
      }

      _.each($scope.categories, function(c) {
        c.Checked = false;
        _.each(userCategories, function(userCategory) {
          if (userCategory.Id === c.Id) {
            c.Checked = true;
          }
        }); // each
      }); // each
      console.log('$scope.userCategories=', $scope.userCategories);
      console.log('$scope.categories=', $scope.categories);
    }) // success
    .error(function(msg, status) {
      console.log(msg);
      console.log('Could not get user machines');
    });
  };

  // TODO: this could be transformed into a filter
  function formatDate(d) {
    //console.log('formatDate: d: ', d);
    var mm = (d.getMonth() + 1);
    if (mm < 10) {
      mm = '0' + mm;
    }
    var dd = d.getDate();
    if (dd < 10) {
      dd = '0' + dd;
    }
    return d.getFullYear() + '-' + mm + '-' + dd;
  }

  $scope.getAvailableMemberships = function() {
    $http({
      method: 'GET',
      url: '/api/memberships',
      params: {
        location: $cookies.get('location'),
        ac: new Date().getTime()
      }
    })
    .success(function(data) {

      $scope.memberships = data;
      $('.adm-user-membership-date').each(function() {
        var eachUserMembershipId = $(this).attr('data-user-membership-id');
        eachUserMembershipId = parseInt(eachUserMembershipId);
        $(this).pickadate({
          format: 'yyyy-mm-dd'
        });
      });

    })
    .error(function(data, status) {
      console.log('Could not get memberships');
      console.log('Data: ' + data);
      console.log('Status code: ' + status);
    });
  };

  $scope.getUserMemberships = function() {
    $http({
      method: 'GET',
      url: '/api/users/' + $scope.user.Id + '/memberships',
      params: {
        ac: new Date().getTime(),
        location: $cookies.get('location')
      }
    })
    .success(function(data) {
      $scope.userMemberships = _.map(data, function(userMembership) {
        
        userMembership.StartDate = 
          new Date(Date.parse(userMembership.StartDate));
        
        userMembership.EndDate = new Date(Date.parse(userMembership.EndDate));
        if (userMembership.TerminationDate && moment(userMembership.TerminationDate).unix() > 0) {
          userMembership.TerminationDateFormatted = moment(userMembership.TerminationDate).format('YYYY-MM-DD');
        }
        var today = new Date();
        
        if (userMembership.StartDate <= today && 
          today <= userMembership.EndDate) {
          
          userMembership.Active = true;

        } else if (today > userMembership.EndDate) {
          if (!userMembership.AutoExtend) {
            userMembership.Cancelled = true;
          } else {
            userMembership.Inactive = true;
          }
        } else {
          userMembership.Inactive = true;
        }

        userMembership.StartDate = formatDate(userMembership.StartDate);
        userMembership.EndDate = formatDate(userMembership.EndDate);
        return userMembership;
      });

      $scope.userMemberships = _.sortBy($scope.userMemberships, function(umb) {
        var t = moment(umb.TerminationDate).unix();

        console.log(t);

        if (t < 0) {
          return -1e20;
        }

        return -t;
      });

      $scope.getAvailableMemberships();
    })
    .error(function(data, status) {
      console.log('Could not get user memberships');
      console.log('Data: ' + data);
      console.log('Status code: ' + status);
    });
  };

  // Adds user membership to the database and updates the UI
  $scope.addUserMembership = function() {
    var startDate = $('#adm-add-user-membership-start-date').val();
    console.log(startDate);
    if (!startDate) {
      toastr.error('Please select a Start Date');
      return;
    }

    var selectedMembershipId = $('#user-select-membership').val();
    if (!selectedMembershipId) {
      toastr.error('Please select a Membership');
      return;
    }
    $http({
      method: 'POST',
      url: '/api/users/' + $scope.user.Id + '/memberships',
      data: {
        startDate: startDate,
        membershipId: selectedMembershipId
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Membership created');
      $scope.getUserMemberships();
    })
    .error(function(text) {
      toastr.error(text);
    });
  };

  $scope.cancel = function() {
    if (confirm('All changes will be discarded, click ok to continue.')) {
      $location.path('/users');
    }
  };

  $scope.deleteUserMembershipPrompt = function(userMembershipId) {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' +
      token + '</span> to delete',
      placeholder: 'Token',
      callback: function(value) {
        if (value) {
          if (value === token) {
            $scope.deleteUserMembership(userMembershipId);
          } else {
            toastr.error('Wrong token');
          }
        } else if (value !== false) {
          toastr.error('No token');
        }
      } // callback
    });
  };

  $scope.updateUserMembershipPrompt = function(userMembershipId) {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' +
      token + '</span> to update',
      placeholder: 'Token',
      callback: function(value) {
        if (value) {
          if (value === token) {
            $scope.updateUserMembership(userMembershipId);
          } else {
            toastr.error('Wrong token');
          }
        } else if (value !== false) {
          toastr.error('No token');
        }
      } // callback
    });
  };


  $scope.updateUserMembership = function(userMembershipId) {
    var userMembership;
    _.each($scope.userMemberships, function(um) {
      if (um.Id && um.Id === userMembershipId) {
        userMembership = um;
      }
    });
    if (userMembership) {
      $http({
        method: 'PUT',
        url: '/api/users/' + $scope.user.Id + '/memberships/' + userMembershipId,
        headers: {'Content-Type': 'application/json' },
        data: userMembership,
        transformRequest: function(data) {
          var transformed = _.extend({}, data);
          var startDate = $('.adm-user-membership-start-date[data-user-membership-id=' + userMembershipId + ']').val();
          if (startDate) {
            transformed.StartDate = startDate;
          }
          var terminationDate = $('.adm-user-membership-termination-date[data-user-membership-id=' + userMembershipId + ']').val();
          if (terminationDate) {
            transformed.TerminationDate = terminationDate;
          } else {
            delete transformed.TerminationDate;
          }
          delete transformed.Invoice;
          return JSON.stringify(transformed);
        },
      })
      .success(function() {
        toastr.success('Membership updated.');
      })
      .error(function(status) {
        toastr.error(status);
      });
    } else {
      toastr.error('Fatal error.');
    }
  };

  $scope.deleteUserMembership = function(userMembershipId) {
    console.log('Delete user membership ID: ' + userMembershipId);
    $http({
      method: 'DELETE',
      url: '/api/users/' + $scope.user.Id + '/memberships/' + userMembershipId,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Membership deleted.');
      $scope.getUserMemberships();
    })
    .error(function(status) {
      toastr.error(status);
    });
  };

  $scope.setTerminationDate = function(userMembershipId) {
    _.each($scope.userMemberships, function(um) {
      if (um.Id && um.Id === userMembershipId) {
        um.TerminationDateFormatted = moment().format('YYYY-MM-DD');
      }
    });
  };

  $scope.unsetTerminationDate = function(userMembershipId) {
    _.each($scope.userMemberships, function(um) {
      if (um.Id && um.Id === userMembershipId) {
        um.TerminationDate = moment('0001-01-01').toDate();
        um.TerminationDateFormatted = null;
        $('.adm-user-membership-termination-date[data-user-membership-id=' + userMembershipId + ']').val('');
      }
    });
  };

  $scope.saveUser = function() {
    if ($scope.user.UserRole === 'admin') {
      $scope.updateUser();
    } else {
      $scope.updateUserMachinePermissions($scope.updateUser);
    }
  };

  $scope.updateUser = function(callback) {
    if (!$scope.user.Email) {
      toastr.info('Please enter an E-Mail');
    }
    if (!$scope.user.Username) {
      toastr.info('Please enter a username');
    }

    $http({
      method: 'PUT',
      url: '/api/users/' + $scope.user.Id,
      headers: {'Content-Type': 'application/json' },
      data: {
        User: $scope.user
      },
      params: {
        location: $cookies.get('location'),
        ac: new Date().getTime()
      },
      transformRequest: function(data) {
        var transformed = {
          User: _.extend({}, data.User)
        };
        _.each(['ClientId', 'FastbillTemplateId' ,'VatRate'], function(field) {
          transformed.User[field] = parseInt(transformed.User[field]);
          if (_.isNaN(transformed.User[field])) {
            transformed.User[field] = 0;
          }
        });
        return JSON.stringify(transformed);
      }
    })
    .success(function() {
      if (callback) {
        callback();
      }
      toastr.success('User updated');
    })
    .error(function(data) {
      if (data === 'duplicateEntry') {
        toastr.error('Duplicate entry error. Make sure that fields like user name and email are unique.');
      } else {
        toastr.error('Error while trying to save changes');
      }
    });
  };

  $scope.setArchived = function(archived) {
    if (archived) {
      $scope.userLocation.UserRole = 'archived';
    } else {
      $scope.userLocation.UserRole = 'member';
    }
    $scope.updateUserLocation($scope.userLocation, function() {
      window.location.reload();
    });
  }; // archive()

  $scope.addUserLocation = function() {
    var locationId = $('.add-user-location select').val();
    if (!locationId) {
      toastr.error('Please select a location to add');
      return;
    }

    $http({
      method: 'POST',
      url: '/api/users/' + $scope.user.Id + '/locations/' + locationId,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('User location created');
      $scope.loadUserLocations();
    })
    .error(function() {
      toastr.error('Error while trying to create new User Location');
    });
  };

  $scope.updateNfcId = function() {
    console.log('updateNfcId()');

    $http({
      method: 'POST',
      url: '/api/users/' + $scope.user.Id + '/nfc_id',
      params: {
        nfcId: $scope.nfcId,
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('NFC ID successfully updated');
    })
    .error(function() {
      toastr.error('Error while trying to update NFC ID');
    });
  };

  $scope.updateUserLocation = function(userLocation, cb) {
    $http({
      method: 'PUT',
      url: '/api/users/' + $scope.user.Id + '/locations/' + userLocation.LocationId,
      headers: {'Content-Type': 'application/json' },
      data: userLocation,
      transformRequest: function(data) {
        return JSON.stringify(data);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.info('Successfully updated user location role.');
      if (cb) {
        cb();
      }
    })
    .error(function() {
      toastr.error('Error while trying to update user location role');
      $scope.loadUserLocations();
    });
  };

  $scope.updateUserMachinePermissions = function(callback) {
    
    // Do not update machine permissions for an admin user
    if ($scope.user.UserRole === 'admin' ||
        $scope.user.UserRole === 'superadmin') {
      if (callback) {
        callback();
      }
      return;
    }

    var permissions = [];
    _.each($scope.availableMachines, function(machine) {
      if (machine.Checked) {
        var p = {MachineId: machine.Id};
        permissions.push(p);
      }
    });

    $http({
      method: 'PUT',
      url: '/api/users/' + $scope.user.Id + '/permissions',
      headers: {'Content-Type': 'application/json' },
      data: permissions,
      transformRequest: function(data) {
        return JSON.stringify(data);
      },
      params: {
        location: $cookies.get('location'),
        ac: new Date().getTime()
      }
    })
    .success(function() {
      if (callback) {
        callback();
      }
    })
    .error(function() {
      toastr.error('Error while trying to update permissions');
    });
  };

  $scope.updatePassword = function() {

    // Check user entered password
    var minPassLength = 3;

    // If there is password at all
    if (!$scope.userPassword || $scope.userPassword === '') {
      toastr.error('No password');
      return;
    }

    // If it is long enough
    if ($scope.userPassword.length < minPassLength) {
      toastr.error('Password too short');
      return;
    }

    // If it matches the repeated password
    if ($scope.userPassword !== $scope.userPasswordRepeat) {
      toastr.error('Passwords do not match');
      return;
    }

    $http({
      method: 'POST', // TODO: change this to PUT as it is an update operation?
      url: '/api/users/' + $scope.user.Id + '/password',
      params: {
        password: $scope.userPassword,
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Password successfully updated');
    })
    .error(function() {
      toastr.error('Error while trying to update password');
    });
  };

  $scope.updateUserRoles = function(value) {
    console.log('value:',value);
  };

  $scope.updateAdminStatus = function() {
    if ($scope.user.UserRole === 'admin' ||
        $scope.user.UserRole === 'superadmin') {
      _.each($scope.availableMachines, function(machine){
        machine.Checked = true;
        machine.Disabled = true;
      });
    } else {
      console.log('uncheck avail machines');
      _.each($scope.availableMachines, function(machine){
        console.log(machine.Id);
        machine.Checked = false;
        machine.Disabled = false;
      });
      $scope.loadUserMachinePermissions();
    }
  };

  // FastBill variables
  var FASTBILL_ACTION_GET_CUSTOMER_NUMBER = '1';
  var FASTBILL_ACTION_LOAD_FROM = '2';
  var FASTBILL_ACTION_UPDATE = '3';
  $scope.fastBillAction = FASTBILL_ACTION_GET_CUSTOMER_NUMBER;

  $scope.readFastBillCustomerNumber = function(data) {
    $scope.user.ClientId = parseInt(data.Customers[0].CUSTOMER_NUMBER);
    $scope.saveUser();
  };

  $scope.readFastBillUserData = function(data) {
    var customer = data.Customers[0];

    $scope.user.FirstName = customer.FIRST_NAME;
    $scope.user.LastName = customer.LAST_NAME;
    $scope.user.InvoiceAddr = customer.ADDRESS;
    $scope.user.ZipCode = customer.ZIPCODE;
    $scope.user.City = customer.CITY;
    $scope.user.CountryCode = customer.COUNTRY_CODE;
    $scope.user.Phone = customer.PHONE;
    $scope.user.ClientId = parseInt(customer.CUSTOMER_NUMBER);

    $scope.saveUser();
  };

  // Sync user data with FastBill account
  $scope.syncWithFastBill = function() {

    if (_.size($scope.user.CountryCode) === 0) {
      toastr.error('Please select a Country Code');
      return;
    }

    console.log($scope.fastBillAction);
    console.log(FASTBILL_ACTION_LOAD_FROM);

    // Check what action the user wants to make
    if (parseInt($scope.fastBillAction) === parseInt(FASTBILL_ACTION_GET_CUSTOMER_NUMBER)) {
      console.log('get customer nr');
      $scope.getFastBillUserByEmail({
        onSuccess: $scope.readFastBillCustomerNumber, 
        onFailure: $scope.createFastBillUser
      });
    }

    if (parseInt($scope.fastBillAction) === parseInt(FASTBILL_ACTION_LOAD_FROM)) {
      console.log('load fb data');
      $scope.getFastBillUserByEmail({
        onSuccess: $scope.readFastBillUserData, 
        onFailure: $scope.createFastBillUser
      });
    }

    if (parseInt($scope.fastBillAction) === parseInt(FASTBILL_ACTION_UPDATE)) {
      console.log('update fb data');
      $scope.updateFastBillUser({
        onSuccess: function(){},
        onFailure: function(){}
      });
    }

  };

  $scope.updateFastBillUser = function(instructions) {

    // Attempt to get fastbill customer ID
    $scope.getFastBillUserByEmail({
      
      // Got the id, use it to update the user
      onSuccess: function(data) {
        var customer = data.Customers[0];

        $http({
          method: 'PUT',
          url: '/api/fastbill/customer/' + customer.CUSTOMER_ID,
          params: {
            firstname: $scope.user.FirstName,
            lastname: $scope.user.LastName,
            //email: 
            address: $scope.user.InvoiceAddr,
            city: $scope.user.City,
            countrycode: $scope.user.CountryCode,
            zipcode: $scope.user.ZipCode,
            phone: $scope.user.Phone,
            //organization:
            ac: new Date().getTime()
          }
        })
        .success(function(data) {
          console.log(data);
          toastr.success('FastBill customer updated.');
          instructions.onSuccess();
        })
        .error(function(data) {
          console.log(data);
          toastr.error('Failed to update FastBill customer.');
          instructions.onFailure();
        });
      },

      // Failed to get FastBill customer id
      onFailure: function() {
        console.log('Failed to get FastBill customer by email.');
        toastr.error('Failed to update FastBill user.');
      }
    });

  };

  $scope.getFastBillUserByEmail = function(instructions) {
    
    if ($scope.user.Email === '') {
      toastr.error('No email address');
      return;
    }

    $http({
      method: 'GET',
      url: '/api/fastbill/customer',
      params: {
        term: $scope.user.Email,
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      if (data.Customers && data.Customers.length) {
        if (data.Customers.length > 1) {
          toastr.warning('More than one FastBill user with the same email address exists.');
        }
        instructions.onSuccess(data);
      } else {
        instructions.onFailure();
      }
    })
    .error(function() {
      toastr.error('An error happened while getting user\'s FastBill customer data.');
    });

  };

  $scope.getFastBillUserByCustomerId = function(customerId, instructions) {
    $http({
      method: 'GET',
      url: '/api/fastbill/customer',
      params: {
        customerid: customerId,
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      if (data.Customers.length) {
        instructions.onSuccess(data);
      } else {
        instructions.onFailure();
      }
    })
    .error(function() {
      toastr.error('An error happened while getting user\'s FastBill customer data.');
    });
  };

  $scope.createFastBillUser = function() {
    toastr.info('Creating new FastBill user.');

    $http({
      method: 'POST',
      url: '/api/fastbill/customer',
      params: {
        firstname: $scope.user.FirstName,
        lastname: $scope.user.LastName,
        email: $scope.user.Email,
        address: $scope.user.InvoiceAddr,
        city: $scope.user.City,
        countrycode: $scope.user.CountryCode,
        zipcode: $scope.user.ZipCode,
        phone: $scope.user.Phone,
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      console.log('Fastbill customer created.');
      console.log(data);
      $scope.getFastBillUserByCustomerId(data.CUSTOMER_ID, {
        onSuccess: $scope.readFastBillCustomerNumber, 
        onFailure: function() {
          toastr.error('Failed to get newly created FastBill customer.');
        }
      });
    })
    .error(function() {
      toastr.error('Failed to create FastBill customer.');
    });
  };

  $scope.hideUserMembershipHistory = function() {
    $scope.userMembershipsHistoryVisible = false;
  };

  $scope.showUserMembershipHistory = function() {
    $scope.userMembershipsHistoryVisible = true;
  };

}]); // app.controller

app.filter('categoriesFilter', function() {
  return function(categories, scope) {
    return _.filter(categories, function(c) {
      return !c.Archived;
    });
  };
});

})(); // closure
