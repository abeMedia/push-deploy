'use strict';

angular.module('app.projects')

.config(function($stateProvider,$httpProvider){
    $stateProvider
    .state('projects.edit',{
      title:'Add/Edit Project',
      url:'/:id/edit',
      templateUrl:'projects/projects.edit.html',
      controller:'ProjectEdit'
    });
})

.controller('ProjectEdit',function($scope,$state,$rootScope,$stateParams,$http,$websocket,popupService,FlashService,Projects){
  if($stateParams.id) {
    $scope.item=Projects.get({id:$stateParams.id});
  }
  else {
    $scope.item=new Projects();
    $scope.item.deploy = [{}];
  }
  
  $scope.removeDeploy = function (i) {
    $scope.item.deploy.splice(i, 1);
  };
  $scope.addDeploy = function () {
    $scope.item.deploy.push({});
  };
  
  $scope.$parent.actions = {
    save: function(){
      if($stateParams.id)
        $scope.item.$update(function(){ $state.go('projects'); });
      else
        $scope.item.$save(function(){ $state.go('projects'); });
    },
    cancel:function(){
      $state.go('projects.view', {id:$scope.item.id});
    }
  };
})

