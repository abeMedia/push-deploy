'use strict';

angular.module('app.projects')

.config(function($stateProvider,$httpProvider){
    $stateProvider
    .state('projects.list',{
      title:'Projects',
      search:true,
      url:'/',
      templateUrl:'projects/projects.list.html',
      controller:'Projects'
    });
})

.controller('Projects',function($scope,$state,$rootScope,$websocket,popupService,Projects){
  $scope.items=Projects.query();
  $scope.$parent.actions={
    new:function(project){
      $state.go('projects.edit');
    }
  };
  $scope.delete=function(project){
    if(popupService.showPopup('Really delete this?')){
      project.$delete();
    }
  };
  $scope.status = statusHelper;
  $rootScope.websocket = $websocket('ws://' + window.location.host + '/api/statuses/' + $rootScope.globals.user.id).onMessage(function(message) {
    var s = message.data.split(":");
    $scope.items[s[0]].status = s[1];
  });
});