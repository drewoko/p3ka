var pekaApp = angular.module('pekaApp', ['ngRoute', 'infinite-scroll']);

pekaApp.controller('aboutController', function($scope, $http) {
});

pekaApp.controller('usersController', function($scope, $http) {
  $http.get('/api/top').then(function(response) {
    if(response.data != null) {
      $scope.users = response.data;
    }
  });
});

pekaApp.controller('userController', function($scope, $http, $routeParams) {

  $scope.currentStart = 0;
  $scope.images = [];

  $scope.showFull = function(self) {
    $('#hoverer img').attr('src', self);
    $('#hoverer').show();
  };

  $('#hoverer').on('click', function(e) {
    if(e.target.localName == "img") {
      $('#hover-img').attr('src', getNextImage($('#hover-img').attr('src'), $scope.images));
    } else {
      $('#hoverer').hide();
    }
  });

  $scope.nextPage = function () {

    $http.get('/api/user?start='+($scope.currentStart)+'&user='+$routeParams.user).then(function(response) {
      if(response.data != null) {
        $scope.currentStart += response.data.length;
        angular.forEach(response.data, function(val) {
          $scope.images.push(val);
        });
      }
    });
  };

  $scope.showFull = function(self) {
    $('#hoverer img').attr('src', self);
    $('#hoverer').show();
  };

  $scope.nextPage();
});

pekaApp.controller('randController', function($scope, $http) {

  $http.get('/api/random').then(function(response) {
    if(response.data != null) {
      $scope.images = response.data;
    }
  });

  $scope.showFull = function(self) {
    $('#hoverer img').attr('src', self);
    $('#hoverer').show();
  };

  $('#hoverer').on('click', function(e) {
    if(e.target.localName == "img") {
      $('#hover-img').attr('src', getNextImage($('#hover-img').attr('src'), $scope.images));
    } else {
      $('#hoverer').hide();
    }
  });
});

pekaApp.controller('mainController', function($scope, $http) {

  $scope.currentStart = 0;
  $scope.images = [];
  $scope.scrollBusy = false;

  $('#hoverer').on('click', function(e) {
    if(e.target.localName == "img") {
      $('#hover-img').attr('src', getNextImage($('#hover-img').attr('src'), $scope.images));
    } else {
      $('#hoverer').hide();
    }
  });

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

  $scope.showFull = function(self) {
    $('#hoverer img').attr('src', self);
    $('#hoverer').show();
  };

  $scope.nextPage();
});

pekaApp.config(function($routeProvider) {
  $routeProvider
    .when('/', {
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
    });
});

function getNextImage(currentSrc, all_images) {

  for(var key in all_images) {

    if(all_images[key].url == currentSrc) {
      return all_images[parseInt(key)+1].url;
    }
  }
  return currentSrc;
}