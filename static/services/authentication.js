'use strict';

angular.module('app')

.factory('AuthenticationService', function ($http, $cookieStore, $rootScope, FlashService) {
  var service = {};

  service.Login = function (username, password, callback) {
    $http.post('/api/login', { username: username, password: password })
    .success(function (user) {
      user.auth = Base64.encode(username + ':' + password);
      $rootScope.globals = {
        user: user
      };
  
      $http.defaults.headers.common['Authorization'] = 'Basic ' + user.auth; // jshint ignore:line
      $cookieStore.put('globals', $rootScope.globals);
      callback(user);
    })
    
    .error(function (user) {
      FlashService("error", "Invalid username or password.")
      console.log(user);
    });
  };
  
  service.Logout = function () {
    $rootScope.globals = {};
    $cookieStore.remove('globals');
    $http.defaults.headers.common.Authorization = 'Basic';
  };

  return service;
})

// http interceptor to sign out if 401 is returned
.config(function ($provide, $httpProvider) {
  $provide.factory('httpInterceptor', function ($q, $window) {
    return {
      response: function (response) {
        return response || $q.when(response);
      },
      responseError: function (rejection) {
        if(rejection.status === 401 && $window.location.hash != '#/login') {
          $window.location.hash = '#/logout';
        }
        return $q.reject(rejection);
      }
    };
  });
  $httpProvider.interceptors.push('httpInterceptor');
});

// Base64 encoding used by AuthenticationService
var Base64 = {

  keyStr: 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=',

  encode: function (input) {
    var output = "";
    var chr1, chr2, chr3 = "";
    var enc1, enc2, enc3, enc4 = "";
    var i = 0;

    do {
      chr1 = input.charCodeAt(i++);
      chr2 = input.charCodeAt(i++);
      chr3 = input.charCodeAt(i++);

      enc1 = chr1 >> 2;
      enc2 = ((chr1 & 3) << 4) | (chr2 >> 4);
      enc3 = ((chr2 & 15) << 2) | (chr3 >> 6);
      enc4 = chr3 & 63;

      if (isNaN(chr2)) {
        enc3 = enc4 = 64;
      } else if (isNaN(chr3)) {
        enc4 = 64;
      }

      output = output +
        this.keyStr.charAt(enc1) +
        this.keyStr.charAt(enc2) +
        this.keyStr.charAt(enc3) +
        this.keyStr.charAt(enc4);
      chr1 = chr2 = chr3 = "";
      enc1 = enc2 = enc3 = enc4 = "";
    } while (i < input.length);

    return output;
  },

  decode: function (input) {
    var output = "";
    var chr1, chr2, chr3 = "";
    var enc1, enc2, enc3, enc4 = "";
    var i = 0;

    // remove all characters that are not A-Z, a-z, 0-9, +, /, or =
    var base64test = /[^A-Za-z0-9\+\/\=]/g;
    if (base64test.exec(input)) {
      window.alert("There were invalid base64 characters in the input text.\n" +
        "Valid base64 characters are A-Z, a-z, 0-9, '+', '/',and '='\n" +
        "Expect errors in decoding.");
    }
    input = input.replace(/[^A-Za-z0-9\+\/\=]/g, "");

    do {
      enc1 = this.keyStr.indexOf(input.charAt(i++));
      enc2 = this.keyStr.indexOf(input.charAt(i++));
      enc3 = this.keyStr.indexOf(input.charAt(i++));
      enc4 = this.keyStr.indexOf(input.charAt(i++));

      chr1 = (enc1 << 2) | (enc2 >> 4);
      chr2 = ((enc2 & 15) << 4) | (enc3 >> 2);
      chr3 = ((enc3 & 3) << 6) | enc4;

      output = output + String.fromCharCode(chr1);

      if (enc3 != 64) {
        output = output + String.fromCharCode(chr2);
      }
      if (enc4 != 64) {
        output = output + String.fromCharCode(chr3);
      }

      chr1 = chr2 = chr3 = "";
      enc1 = enc2 = enc3 = enc4 = "";

    } while (i < input.length);

    return output;
  }
};
