(function(){

'use strict';

angular.module('fabsmith.signup.form', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/form', {
    templateUrl: 'ng-modules/form/form.html',
    controller: 'FormCtrl'
  });
}])

.controller('FormCtrl',
 ['$scope', '$location', '$http', '$q', '$routeParams',
 function($scope, $location, $http, $q, $routeParams) {

  // Regular expression for email spec : RFC 5322
  $scope.emailRegExp = /^[-a-z0-9~!$%^*_=+}{\'?]+(\.[-a-z0-9~!$%^*_=+}{\'?]+)*@([a-z0-9_][-a-z0-9_]*(\.[-a-z0-9_]+)*\.(aero|arpa|biz|com|coop|edu|gov|info|int|mil|museum|name|net|org|pro|travel|mobi|[a-z][a-z])|([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))(:[0-9]{1,5})?/i;
  $scope.minUsernameAndPasswordLength = 3;

  var countryCodesJSON = '[{"Name":"Afghanistan","Code":"AF"},{"Name":"Åland Islands","Code":"AX"},{"Name":"Albania","Code":"AL"},{"Name":"Algeria","Code":"DZ"},{"Name":"American Samoa","Code":"AS"},{"Name":"Andorra","Code":"AD"},{"Name":"Angola","Code":"AO"},{"Name":"Anguilla","Code":"AI"},{"Name":"Antarctica","Code":"AQ"},{"Name":"Antigua and Barbuda","Code":"AG"},{"Name":"Argentina","Code":"AR"},{"Name":"Armenia","Code":"AM"},{"Name":"Aruba","Code":"AW"},{"Name":"Australia","Code":"AU"},{"Name":"Austria","Code":"AT"},{"Name":"Azerbaijan","Code":"AZ"},{"Name":"Bahamas","Code":"BS"},{"Name":"Bahrain","Code":"BH"},{"Name":"Bangladesh","Code":"BD"},{"Name":"Barbados","Code":"BB"},{"Name":"Belarus","Code":"BY"},{"Name":"Belgium","Code":"BE"},{"Name":"Belize","Code":"BZ"},{"Name":"Benin","Code":"BJ"},{"Name":"Bermuda","Code":"BM"},{"Name":"Bhutan","Code":"BT"},{"Name":"Bolivia, Plurinational State of","Code":"BO"},{"Name":"Bonaire, Sint Eustatius and Saba","Code":"BQ"},{"Name":"Bosnia and Herzegovina","Code":"BA"},{"Name":"Botswana","Code":"BW"},{"Name":"Bouvet Island","Code":"BV"},{"Name":"Brazil","Code":"BR"},{"Name":"British Indian Ocean Territory","Code":"IO"},{"Name":"Brunei Darussalam","Code":"BN"},{"Name":"Bulgaria","Code":"BG"},{"Name":"Burkina Faso","Code":"BF"},{"Name":"Burundi","Code":"BI"},{"Name":"Cambodia","Code":"KH"},{"Name":"Cameroon","Code":"CM"},{"Name":"Canada","Code":"CA"},{"Name":"Cape Verde","Code":"CV"},{"Name":"Cayman Islands","Code":"KY"},{"Name":"Central African Republic","Code":"CF"},{"Name":"Chad","Code":"TD"},{"Name":"Chile","Code":"CL"},{"Name":"China","Code":"CN"},{"Name":"Christmas Island","Code":"CX"},{"Name":"Cocos (Keeling) Islands","Code":"CC"},{"Name":"Colombia","Code":"CO"},{"Name":"Comoros","Code":"KM"},{"Name":"Congo","Code":"CG"},{"Name":"Congo, the Democratic Republic of the","Code":"CD"},{"Name":"Cook Islands","Code":"CK"},{"Name":"Costa Rica","Code":"CR"},{"Name":"Côte d\'Ivoire","Code":"CI"},{"Name":"Croatia","Code":"HR"},{"Name":"Cuba","Code":"CU"},{"Name":"Curaçao","Code":"CW"},{"Name":"Cyprus","Code":"CY"},{"Name":"Czech Republic","Code":"CZ"},{"Name":"Denmark","Code":"DK"},{"Name":"Djibouti","Code":"DJ"},{"Name":"Dominica","Code":"DM"},{"Name":"Dominican Republic","Code":"DO"},{"Name":"Ecuador","Code":"EC"},{"Name":"Egypt","Code":"EG"},{"Name":"El Salvador","Code":"SV"},{"Name":"Equatorial Guinea","Code":"GQ"},{"Name":"Eritrea","Code":"ER"},{"Name":"Estonia","Code":"EE"},{"Name":"Ethiopia","Code":"ET"},{"Name":"Falkland Islands (Malvinas)","Code":"FK"},{"Name":"Faroe Islands","Code":"FO"},{"Name":"Fiji","Code":"FJ"},{"Name":"Finland","Code":"FI"},{"Name":"France","Code":"FR"},{"Name":"French Guiana","Code":"GF"},{"Name":"French Polynesia","Code":"PF"},{"Name":"French Southern Territories","Code":"TF"},{"Name":"Gabon","Code":"GA"},{"Name":"Gambia","Code":"GM"},{"Name":"Georgia","Code":"GE"},{"Name":"Germany","Code":"DE"},{"Name":"Ghana","Code":"GH"},{"Name":"Gibraltar","Code":"GI"},{"Name":"Greece","Code":"GR"},{"Name":"Greenland","Code":"GL"},{"Name":"Grenada","Code":"GD"},{"Name":"Guadeloupe","Code":"GP"},{"Name":"Guam","Code":"GU"},{"Name":"Guatemala","Code":"GT"},{"Name":"Guernsey","Code":"GG"},{"Name":"Guinea","Code":"GN"},{"Name":"Guinea-Bissau","Code":"GW"},{"Name":"Guyana","Code":"GY"},{"Name":"Haiti","Code":"HT"},{"Name":"Heard Island and McDonald Islands","Code":"HM"},{"Name":"Holy See (Vatican City State)","Code":"VA"},{"Name":"Honduras","Code":"HN"},{"Name":"Hong Kong","Code":"HK"},{"Name":"Hungary","Code":"HU"},{"Name":"Iceland","Code":"IS"},{"Name":"India","Code":"IN"},{"Name":"Indonesia","Code":"ID"},{"Name":"Iran, Islamic Republic of","Code":"IR"},{"Name":"Iraq","Code":"IQ"},{"Name":"Ireland","Code":"IE"},{"Name":"Isle of Man","Code":"IM"},{"Name":"Israel","Code":"IL"},{"Name":"Italy","Code":"IT"},{"Name":"Jamaica","Code":"JM"},{"Name":"Japan","Code":"JP"},{"Name":"Jersey","Code":"JE"},{"Name":"Jordan","Code":"JO"},{"Name":"Kazakhstan","Code":"KZ"},{"Name":"Kenya","Code":"KE"},{"Name":"Kiribati","Code":"KI"},{"Name":"Korea, Democratic People\'s Republic of","Code":"KP"},{"Name":"Korea, Republic of","Code":"KR"},{"Name":"Kuwait","Code":"KW"},{"Name":"Kyrgyzstan","Code":"KG"},{"Name":"Lao People\'s Democratic Republic","Code":"LA"},{"Name":"Latvia","Code":"LV"},{"Name":"Lebanon","Code":"LB"},{"Name":"Lesotho","Code":"LS"},{"Name":"Liberia","Code":"LR"},{"Name":"Libya","Code":"LY"},{"Name":"Liechtenstein","Code":"LI"},{"Name":"Lithuania","Code":"LT"},{"Name":"Luxembourg","Code":"LU"},{"Name":"Macao","Code":"MO"},{"Name":"Macedonia, the Former Yugoslav Republic of","Code":"MK"},{"Name":"Madagascar","Code":"MG"},{"Name":"Malawi","Code":"MW"},{"Name":"Malaysia","Code":"MY"},{"Name":"Maldives","Code":"MV"},{"Name":"Mali","Code":"ML"},{"Name":"Malta","Code":"MT"},{"Name":"Marshall Islands","Code":"MH"},{"Name":"Martinique","Code":"MQ"},{"Name":"Mauritania","Code":"MR"},{"Name":"Mauritius","Code":"MU"},{"Name":"Mayotte","Code":"YT"},{"Name":"Mexico","Code":"MX"},{"Name":"Micronesia, Federated States of","Code":"FM"},{"Name":"Moldova, Republic of","Code":"MD"},{"Name":"Monaco","Code":"MC"},{"Name":"Mongolia","Code":"MN"},{"Name":"Montenegro","Code":"ME"},{"Name":"Montserrat","Code":"MS"},{"Name":"Morocco","Code":"MA"},{"Name":"Mozambique","Code":"MZ"},{"Name":"Myanmar","Code":"MM"},{"Name":"Namibia","Code":"NA"},{"Name":"Nauru","Code":"NR"},{"Name":"Nepal","Code":"NP"},{"Name":"Netherlands","Code":"NL"},{"Name":"New Caledonia","Code":"NC"},{"Name":"New Zealand","Code":"NZ"},{"Name":"Nicaragua","Code":"NI"},{"Name":"Niger","Code":"NE"},{"Name":"Nigeria","Code":"NG"},{"Name":"Niue","Code":"NU"},{"Name":"Norfolk Island","Code":"NF"},{"Name":"Northern Mariana Islands","Code":"MP"},{"Name":"Norway","Code":"NO"},{"Name":"Oman","Code":"OM"},{"Name":"Pakistan","Code":"PK"},{"Name":"Palau","Code":"PW"},{"Name":"Palestine, State of","Code":"PS"},{"Name":"Panama","Code":"PA"},{"Name":"Papua New Guinea","Code":"PG"},{"Name":"Paraguay","Code":"PY"},{"Name":"Peru","Code":"PE"},{"Name":"Philippines","Code":"PH"},{"Name":"Pitcairn","Code":"PN"},{"Name":"Poland","Code":"PL"},{"Name":"Portugal","Code":"PT"},{"Name":"Puerto Rico","Code":"PR"},{"Name":"Qatar","Code":"QA"},{"Name":"Réunion","Code":"RE"},{"Name":"Romania","Code":"RO"},{"Name":"Russian Federation","Code":"RU"},{"Name":"Rwanda","Code":"RW"},{"Name":"Saint Barthélemy","Code":"BL"},{"Name":"Saint Helena, Ascension and Tristan da Cunha","Code":"SH"},{"Name":"Saint Kitts and Nevis","Code":"KN"},{"Name":"Saint Lucia","Code":"LC"},{"Name":"Saint Martin (French part)","Code":"MF"},{"Name":"Saint Pierre and Miquelon","Code":"PM"},{"Name":"Saint Vincent and the Grenadines","Code":"VC"},{"Name":"Samoa","Code":"WS"},{"Name":"San Marino","Code":"SM"},{"Name":"Sao Tome and Principe","Code":"ST"},{"Name":"Saudi Arabia","Code":"SA"},{"Name":"Senegal","Code":"SN"},{"Name":"Serbia","Code":"RS"},{"Name":"Seychelles","Code":"SC"},{"Name":"Sierra Leone","Code":"SL"},{"Name":"Singapore","Code":"SG"},{"Name":"Sint Maarten (Dutch part)","Code":"SX"},{"Name":"Slovakia","Code":"SK"},{"Name":"Slovenia","Code":"SI"},{"Name":"Solomon Islands","Code":"SB"},{"Name":"Somalia","Code":"SO"},{"Name":"South Africa","Code":"ZA"},{"Name":"South Georgia and the South Sandwich Islands","Code":"GS"},{"Name":"South Sudan","Code":"SS"},{"Name":"Spain","Code":"ES"},{"Name":"Sri Lanka","Code":"LK"},{"Name":"Sudan","Code":"SD"},{"Name":"Suriname","Code":"SR"},{"Name":"Svalbard and Jan Mayen","Code":"SJ"},{"Name":"Swaziland","Code":"SZ"},{"Name":"Sweden","Code":"SE"},{"Name":"Switzerland","Code":"CH"},{"Name":"Syrian Arab Republic","Code":"SY"},{"Name":"Taiwan, Province of China","Code":"TW"},{"Name":"Tajikistan","Code":"TJ"},{"Name":"Tanzania, United Republic of","Code":"TZ"},{"Name":"Thailand","Code":"TH"},{"Name":"Timor-Leste","Code":"TL"},{"Name":"Togo","Code":"TG"},{"Name":"Tokelau","Code":"TK"},{"Name":"Tonga","Code":"TO"},{"Name":"Trinidad and Tobago","Code":"TT"},{"Name":"Tunisia","Code":"TN"},{"Name":"Turkey","Code":"TR"},{"Name":"Turkmenistan","Code":"TM"},{"Name":"Turks and Caicos Islands","Code":"TC"},{"Name":"Tuvalu","Code":"TV"},{"Name":"Uganda","Code":"UG"},{"Name":"Ukraine","Code":"UA"},{"Name":"United Arab Emirates","Code":"AE"},{"Name":"United Kingdom","Code":"GB"},{"Name":"United States","Code":"US"},{"Name":"United States Minor Outlying Islands","Code":"UM"},{"Name":"Uruguay","Code":"UY"},{"Name":"Uzbekistan","Code":"UZ"},{"Name":"Vanuatu","Code":"VU"},{"Name":"Venezuela, Bolivarian Republic of","Code":"VE"},{"Name":"Viet Nam","Code":"VN"},{"Name":"Virgin Islands, British","Code":"VG"},{"Name":"Virgin Islands, U.S.","Code":"VI"},{"Name":"Wallis and Futuna","Code":"WF"},{"Name":"Western Sahara","Code":"EH"},{"Name":"Yemen","Code":"YE"},{"Name":"Zambia","Code":"ZM"},{"Name":"Zimbabwe","Code":"ZW"}]';
  $scope.countryCodes = JSON.parse(countryCodesJSON);

  $scope.locationId = parseInt($routeParams.location || 1);
  $scope.location = {};

  $scope.getData = function() {
    $q.all([
      $http({
        method: 'GET',
        url: '/api/locations',
        params: {
          ac: new Date().getTime()
        }
      }),
      $http({
        method: 'GET',
        url: 'http://localhost:8080/api/settings/terms_url',
        params: {
          location: $scope.locationId,
          ac: new Date().getTime()
        }
      })
    ])
    .then(function(resp) {
      var locations = resp[0].data;
      var termsUrl = resp[1].data;

      $scope.locations = locations;
      _.each($scope.locations, function(location) {
        if (location.Id === $scope.locationId) {
          $scope.location = location;
        }
      });
      $scope.termsUrl = termsUrl;
    })
    .catch(function() {
      toastr.error('Failed to load data');
    });
  };

  $scope.getData();

  $scope.submitForm = function() {

    // Check if email is empty
    if (!$scope.email) {
      toastr.error('Please enter an e-mail.');
      return;
    }

    // Check if E-Mail is valid
    if (!$scope.emailRegExp.test($scope.email)) {
      toastr.error('Please enter a valid e-mail.');
      return;
    }

    // Check if First Name is empty
    if (!$scope.firstName) {
      toastr.error('Please enter a first name.');
      return;
    }

    // Check if Last Name is empty
    if (!$scope.lastName) {
      toastr.error('Please enter a last name.');
      return;
    }

    // Check if Username is empty
    if (!$scope.username) {
      toastr.error('Please enter an username.');
      return;
    }

    // Check if Username is long enough
    if ($scope.username.length < $scope.minUsernameAndPasswordLength) {
      toastr.error('Please enter an username which is at least 3 characters long.');
      return;
    }

    // Check if phone number is entered
    if (!$scope.phone) {
      toastr.error('Please enter your phone number.');
      return;
    }

    // Check if Password is empty
    if (!$scope.password) {
      toastr.error('Please enter a password.');
      return;
    }

    // Check if Password is long enough
    if ($scope.password.length < $scope.minUsernameAndPasswordLength) {
      toastr.error('Please enter a password which is at least 3 characters long.');
      return;
    }

    // Check if password matches the repeated password
    if ($scope.password !== $scope.passwordRepeat) {
      toastr.error('Passwords do not match.');
      return;
    }

    // Check if invoice address is added
    if (!$scope.invoiceAddr) {
      toastr.error('Please provide your billing address.');
      return;
    }

    // Check if country is selected
    if (!$scope.countryCode) {
      toastr.error('Please select your country.');
      return;
    }

    // Check if ZIP code is entered
    if (!$scope.zipCode) {
      toastr.error('Please enter your ZIP code.');
      return;
    }

    // Check if city is provided
    if (!$scope.city) {
      toastr.error('Please enter your city.');
      return;
    }

    // Check if AGB and Data Protection Agreement is checked
    if (!$scope.agb_dpa_agreed) {
      toastr.error('You have to agree to the AGB and Data Protection Agreement.');
      return;
    }

    $http({
      method: 'POST',
      url: '/api/users/signup?location=' + $scope.locationId,
      data: {
        "User": {
          "Email": $scope.email,
          "Company": $scope.company,
          "FirstName": $scope.firstName,
          "LastName": $scope.lastName,
          "Username": $scope.username,
          "InvoiceAddr": $scope.invoiceAddr,
          "Phone": $scope.phone,
          "CountryCode": $scope.countryCode,
          "City": $scope.city,
          "ZipCode": $scope.zipCode
        },
        "Password": $scope.password,
        // TODO: Remove this: it is client time - we trust only server time.
        "RegistrationDate": Date.now()
      },
      transformRequest: function(data) {
        return JSON.stringify(data);
      },
      headers: {'Content-Type': 'application/json'}
    })
    .success(function() {
      toastr.success("Registration successful !");
      $location.path('/thanks');
    })
    .error(function() {
      toastr.error('Error while trying to register');
    });
  };

  var numberOfSecondsBeforeResting = 30;
  (function(nbSecToReset){
    var idleTime = 0;
    $(document).ready(function () {
      var idleInterval = setInterval(timerIncrement, 1000);

      $(this).mousemove(function (e) {
        idleTime = 0;
      });
      $(this).keypress(function (e) {
        idleTime = 0;
      });
    });

    function timerIncrement() {
      idleTime = idleTime + 1;
      if (idleTime > (nbSecToReset-1)) {
        $scope.email = "";
        $scope.company = "";
        $scope.firstName = "";
        $scope.lastName = "";
        $scope.username = "";
        $scope.password = "";
        $scope.$apply();
        idleTime = 0;
      }
    }
  })(numberOfSecondsBeforeResting);

}]);
})();
