'use strict';

angular.module('app.projects')

.config(function($stateProvider,$httpProvider){
    $stateProvider
    .state('projects.builds',{
      title:'Build History',
      url:'/:project_id/builds',
      templateUrl:'projects/projects.builds.html',
      controller:'ProjectBuilds',
      subMenu: ["projects.view","projects.builds"],
      icon:'clock-o'
    });
})

.controller('ProjectBuilds',function($scope,$state,$rootScope,$stateParams,$http,$websocket,popupService,FlashService,Builds){
  $scope.items=Builds.query({project_id:$stateParams.project_id});
  $scope.status = statusHelper;
});
