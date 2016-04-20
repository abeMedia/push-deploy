'use strict';

angular.module('app.projects', [
  'ui.router',
  'angular-websocket',
])

.config(function($stateProvider,$httpProvider){
    $stateProvider
    .state('projects',{
      url:'/projects',
      templateUrl:'base/dashboard.html',
      redirectTo: 'projects.list',
    })
    .state('projects.list',{
      title:'Projects',
      search:true,
      url:'/',
      templateUrl:'projects/projects.list.html',
      controller:'Projects'
    })
    .state('projects.view',{
      title:'Project Details',
      url:'/:id',
      templateUrl:'projects/projects.view.html',
      controller:'ProjectView'
    })
    .state('projects.edit',{
      title:'Add/Edit Project',
      url:'/:id/edit',
      templateUrl:'projects/projects.edit.html',
      controller:'ProjectEdit'
    });
})

.factory('Projects',function($resource, $state){
  return $resource('/api/projects/:id',{id:'@id'},{
    update: {
      method: 'PUT'
    },
    build:  {
      method: 'POST',
      url: '/api/projects/:id/build',
      interceptor: {
        response: function (data, q) {
          console.log('response in interceptor', data);
          console.log('q', q);
        },
        responseError: function (data) {
          console.log('error in interceptor', data);
        }
      },
    }
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

.controller('ProjectView',function($scope,$state,$rootScope,$stateParams,$http,$websocket,popupService,FlashService,Projects){
  $scope.item=Projects.get({id:$stateParams.id});
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

var statusHelper = {
  icon: function(id) {
    var icons = {
      "-1":"fa fa-times",
      "0":"fa fa-minus",
      "1":"fa fa-check",
      "2":"fa fa-cog fa-spin",
      "3":"fa fa-cloud-upload",
    };
    return icons[id] + " " + statusHelper.class(id);
  },
  text: function(id) {
    var strings = {
      "-1":"Error!",
      "0":"No builds yet",
      "1":"All dandy!",
      "2":"Building...",
      "3":"Deploying...",
    };
    return strings[id];
  },
  class: function(id) {
    var classes = {
      "-1":"text-danger",
      "0":"text-warning",
      "1":"text-success",
      "2":"text-info",
      "3":"text-info",
    };
    return classes[id];
  }
};