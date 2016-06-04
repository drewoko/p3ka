var pekaApp = angular.module('pekaApp', ['ngRoute']);

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

  var lastScrollHeigth = document.body.scrollHeight;

  $(window).scroll(function() {

    var totalScroll = document.body.scrollHeight-$(window).height();
    var scrollTop = $(window).scrollTop();

    if(totalScroll/2 < scrollTop && lastScrollHeigth < document.body.scrollHeight) {
      $scope.nextPage();
      lastScrollHeigth = document.body.scrollHeight;
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

  $scope.showFull = function(self, next) {
    $('#hoverer img').attr({
      'src': self,
      'data-next': next
    });
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

  $scope.showFull = function(self, next) {
    $('#hoverer img').attr({
      'src': self,
      'data-next': next
    });
    $('#hoverer').show();
  };
});

pekaApp.controller('mainController', function($scope, $http) {

  $scope.currentStart = 0;
  $scope.images = [];

  var lastScrollHeigth = document.body.scrollHeight;

  $(window).scroll(function() {

    var totalScroll = document.body.scrollHeight-$(window).height();
    var scrollTop = $(window).scrollTop();

    if(totalScroll/2 < scrollTop && lastScrollHeigth < document.body.scrollHeight) {
      $scope.nextPage();
      lastScrollHeigth = document.body.scrollHeight;
    }
  });

  $scope.nextPage = function () {

    $http.get('/api/last?start='+($scope.currentStart)).then(function(response) {
      if(response.data != null) {
        $scope.currentStart += response.data.length;
        angular.forEach(response.data, function(val) {
          $scope.images.push(val);
        });
      }
    });
  };

  $scope.showFull = function(self, next) {
    $('#hoverer img').attr({
      'src': self,
      'data-next': next
    });
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

$(function() {
  $('#hoverer').on('click', function(e) {
      $('#hoverer').hide();
  });
});