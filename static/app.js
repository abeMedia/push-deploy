'use strict';

// Declare app level module which depends on views, and components
var app = angular.module('app', [
  'ui.router',
  'ui.bootstrap',
  
  // for page-loading animations
  'angular-loading-bar',
  'ngAnimate', 
  
  // for auth
  'ngCookies', 
  
  // for api
  'ngResource', 
  
  'app.login', 
  'app.projects',
  'app.users', 
])

.config(['cfpLoadingBarProvider', function(cfpLoadingBarProvider) {
  cfpLoadingBarProvider.includeSpinner = false;
  cfpLoadingBarProvider.latencyThreshold = 200;
}])

.run(['$rootScope', '$state', '$cookieStore', '$http', function($rootScope, $state, $cookieStore, $http) {
    // keep user logged in after page refresh
    $rootScope.globals = $cookieStore.get('globals') || {};
    if ($rootScope.globals.user) {
      $http.defaults.headers.common['Authorization'] = 'Basic ' + $rootScope.globals.user.auth; // jshint ignore:line
    }
    
    if($state.current.abstract) $state.go('projects.list');
    
    
    $rootScope.$on('$stateChangeStart', function(evt, to, params) {
      // close open websocket
      if($rootScope.websocket) $rootScope.websocket.close();
      
      // redirect to login page if not logged in and trying to access a restricted page
      var restrictedPage = ['login', 'register'].indexOf(to.name) === -1;
      var loggedIn = $rootScope.globals.user;
      if (restrictedPage && !loggedIn) {
        evt.preventDefault();
        $state.go('login');
      }
      
      // redirect to sub-state if `redirectTo` is set
      if (to.redirectTo) {
        evt.preventDefault();
        $state.go(to.redirectTo, params);
      }
    });
    
    // add page title and actions to rootscope
    angular.forEach([ '$stateChangeSuccess', '$stateChangeError'], function(event) {
      $rootScope.$on(event, function(event, toState, toParams, fromState, fromParams, error) {
        $rootScope.page = $state.current;
      });
    });
    
    // add loading animation to logo
    $rootScope.$on('cfpLoadingBar:started', function(evt, to, params) {
      angular.element(document.getElementById("logo")).addClass('loading');
    });
    $rootScope.$on('cfpLoadingBar:completed', function(evt, to, params) {
      setTimeout(function(){ angular.element(document.getElementById("logo")).removeClass('loading'); }, 500);
    });
}]);