var pekaApp = angular.module('pekaApp', ['ngRoute', 'infinite-scroll', 'ngLocationUpdate']);

pekaApp.controller('aboutController', function($scope, $http) {
});

pekaApp.controller('usersController', function($scope, $http) {
  $http.get('/api/top').then(function(response) {
    if(response.data != null) {
      $scope.users = response.data;
    }
  });
});

pekaApp.controller('userController', function($scope, $rootScope, $http, $routeParams) {

  $scope.currentStart = 0;
  $scope.images = [];
  $scope.scrollBusy = false;

  $scope.nextPage = function () {
    if($scope.scrollBusy) return;

    $scope.scrollBusy = true;
    $http.get('/api/user?start='+($scope.currentStart)+'&user='+$routeParams.user).then(function(response) {
      if(response.data != null) {
        $scope.currentStart += response.data.length;
        angular.forEach(response.data, function(val) {
          $scope.images.push(val);
        });
        $scope.scrollBusy = false;
      }
    });
  };

  $scope.showFull = function (image, playlist) {
    $rootScope.LightBox.show(image.url, playlist);
  };

  $scope.nextPage();
});

pekaApp.controller('randController', function($scope, $rootScope, $http) {

  $http.get('/api/random').then(function(response) {
    if(response.data != null) {
      $scope.images = response.data;
    }
  });

  $scope.showFull = function (image, playlist) {
    $rootScope.LightBox.show(image.url, playlist);
  };
});

pekaApp.controller('mainController', function($scope, $rootScope, $http, $routeParams) {

  if($routeParams.image != undefined) {
    $rootScope.LightBox.show(atob($routeParams.image), []);
  }

  $scope.currentStart = 0;
  $scope.images = [];
  $scope.scrollBusy = false;

  $scope.showFull = function (image, playlist) {
    $rootScope.LightBox.show(image.url, playlist);
  };

  $scope.nextPage = function () {
    if (this.scrollBusy) return;

    this.scrollBusy = true;

    $http.get('/api/last?start='+($scope.currentStart)).then(function(response) {
      if(response.data != null) {
        $scope.currentStart += response.data.length;
        angular.forEach(response.data, function(val) {
          $scope.images.push(val);
        });
        $scope.scrollBusy = false;
      }
    });
  };

  $scope.nextPage();
});

pekaApp.controller('LightBox', function($rootScope, $scope, $window, $location) {

  this.url = "";
  this.show = false;

  var origUrl;
  var imagePlayList = [];

  $rootScope.LightBox = {};
  $rootScope.LightBox.show = function(url, playlist) {

    origUrl = getOrigUrl();

    $scope.url = url;
    $scope.show = true;
    $location.update_path("/show/"+btoa(url), false);
    imagePlayList = playlist;
  };

  $scope.closeClick = function($event) {
    if($event.target.localName != "img") {
      close();
      $location.update_path(origUrl, false);
    }
  };

  $scope.next = function () {
    this.url = getNextImage($scope.url, imagePlayList);
    $location.update_path("/show/"+btoa(this.url), false);
  };

  function close() {
    $scope.url = "";
    $scope.show = false;
    imagePlayList = [];
  }

  function getOrigUrl() {
    var origUrl = $window.location.hash;
    return origUrl.indexOf("/show") > -1 ? "/" : origUrl.replace('#/', '');
  }
});

pekaApp.config(function($routeProvider) {
  $routeProvider
    .when('/', {
      templateUrl : 'pages/main.html',
      controller  : 'mainController'
    })
    .when('/show/:image', {
      templateUrl : 'pages/main.html',
      controller  : 'mainController'
    })
    .when('/random', {
      templateUrl : 'pages/random.html',
      controller  : 'randController'
    })
    .when('/user/:user/', {
      templateUrl : 'pages/user.html',
      controller  : 'userController'
    })
    .when('/users', {
      templateUrl : 'pages/users.html',
      controller  : 'usersController'
    })
    .when('/about', {
      templateUrl : 'pages/about.html',
      controller  : 'aboutController'
    })
    .otherwise({
      redirectTo: '/'
    });
});

function getNextImage(currentSrc, all_images) {

  for(var key in all_images) {

    if(all_images[key].url == currentSrc) {
      var k = parseInt(key)+1;
      if(all_images[k] != undefined) {
        return all_images[k].url;
      }
    }
  }
  return currentSrc;
}