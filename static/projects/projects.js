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
      redirectTo: 'projects.list'
    });
})

.factory('Projects',function($resource, $state){
  return $resource('/api/projects/:id',{id:'@id'},{
    update: {
      method: 'PUT'
    }
  });
})

.factory('Builds',function($resource, $state){
  return $resource('/api/projects/:project_id/builds/:id',{project_id:'@project_id',id:'@id'});
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