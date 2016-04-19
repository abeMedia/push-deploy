'use strict';

angular.module('app.login', [
  'ui.router',
])

.config(function($stateProvider,$httpProvider){
  $stateProvider
  .state('login',{
    url:'/login',
    templateUrl:'login/login.html',
    controller: function ($scope, $state, AuthenticationService, FlashService) {
      (function initController() {
        // reset login status
        AuthenticationService.Logout();
      })();
      
      $scope.login = function () {
        AuthenticationService.Login($scope.username, $scope.password, function (response) {
          $state.go('projects');
        });
      };
    }
  })
  .state('logout',{
    url:'/logout',
    controller: function ($scope, $state, AuthenticationService, FlashService) {
      AuthenticationService.Logout();
      $state.go('login');
    }
  });
});
