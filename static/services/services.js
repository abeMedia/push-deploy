angular.module('app')

.factory('Users',function($resource){
  return $resource('/api/users/:id',{id:'@id'},{
    update: {
      method: 'PUT'
    }
  });
})

.service('popupService',function($window){
  this.showPopup=function(message){
    return $window.confirm(message);
  };
});