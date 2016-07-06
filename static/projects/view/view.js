'use strict';

angular.module('app.projects')

.config(function($stateProvider,$httpProvider){
    $stateProvider
    .state('projects.view',{
      title:'...',
      url:'/:id',
      templateUrl:'projects/projects.view.html',
      controller:'ProjectView',
      subMenu: ["projects.view({id:item.id})","projects.builds({id:item.id})"],
      icon:'flag'
    });
})

.controller('ProjectView',function($scope,$state,$rootScope,$stateParams,$http,$websocket,popupService,FlashService,Projects){
  $scope.item=Projects.get({id:$stateParams.id}, function() {
      $rootScope.page.title = $scope.item.name;
  });
  $rootScope.websocket = $websocket('ws://' + window.location.host + '/api/status/' + $stateParams.id).onMessage(function(message) {
    $scope.item.status = message.data;
  });
  $scope.webhook_url = window.location.href.split("#")[0] + "hook/" + $stateParams.id;
  
  $scope.$parent.actions = {
    edit:function(){
      $state.go('projects.edit',{id:$scope.item.id});
    },
    delete:function(){
      if(popupService.showPopup('Really delete this?')){
        $scope.item.$delete(function(){
          $state.go('projects');
        });
      }
    }
  };
  
  $scope.status = statusHelper;
});
