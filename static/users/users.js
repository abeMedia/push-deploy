'use strict';

angular.module('app.users', [
  'ui.router',
  'ngPassword',
])

.factory('Users',function($resource){
  return $resource('/api/users/:id',{id:'@id'},{
    update: {
      method: 'PUT'
    }
  });
})

.config(function($stateProvider,$httpProvider){
  $stateProvider
  .state('users',{
    url:'/users',
    templateUrl:'base/dashboard.html',
    redirectTo: 'users.list',
  })
  
  .state('users.list',{
    url:'/',
    title: "Users",
    search:true,
    actions: ["new"],
    templateUrl:'users/users.list.html',
    controller:function($scope,$state,popupService,Users){
      $scope.items=Users.query();
      $scope.delete=function(item){
        if(popupService.showPopup('Really delete this?')){
          item.$delete();
        }
      };
      $scope.$parent.actions = {
        new:function(){
          $state.go('users.edit');
        }
      };
    }
  })
  
  .state('users.edit',{
    url:'/edit/:id',
    title: "Add/Edit User",
    actions: ["save","cancel"],
    templateUrl:'users/users.edit.html',
    controller:function($scope,$state,Users){
      if($state.params.id) {
        $scope.item=Users.get({id:$state.params.id});
      }
      else {
        $scope.item=new Users();
      }
      
      $scope.$parent.actions = {
        save: function(){
          if($state.params.id)
            $scope.item.$update(function(){ $state.go('users'); });
          else
            $scope.item.$save(function(){ $state.go('users'); });
        },
        cancel:function(){
          $state.go('users.list', {id:$scope.item.id});
        },
      };
    }
  })
  
  .state('users.profile',{
    url:'/profile',
    title: "Edit Profile",
    actions: ["save","cancel"],
    templateUrl:'users/users.edit.html',
    controller:function($scope,$state,$cookieStore,$rootScope,Users){
      $scope.item=Users.get({id:$rootScope.globals.user.id});
      
      $scope.$parent.actions = {
        save:function(){ 
          $scope.item.$update(function(user){
            // remove functions from user object
            user = JSON.parse(JSON.stringify(user));
            
            // merge with current user object
            for (var field in user) {
              $rootScope.globals.user[field] = user[field];
            }
            $cookieStore.put('globals', $rootScope.globals);
            
            $state.go('users');
          });
        },
        cancel:function(){ 
          $state.go('users');
        },
      };
    }
  });
});