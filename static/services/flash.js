'use strict';

angular.module('app')

.factory('FlashService', function ($rootScope) {
  $rootScope.$on('$locationChangeStart', function () {
    clearFlashMessage();
  });

  function clearFlashMessage() {
    var flash = $rootScope.flash;
    if (flash) {
      if (!flash.sticky) {
        delete $rootScope.flash;
      } else {
        // only keep for a single location change
        flash.sticky = false;
      }
    }
  }

  return function (type, message, sticky) {
    $rootScope.flash = {
      message: message,
      type: (type == 'error' ? 'danger' : type),
      sticky: sticky
    };
  };
});
