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
  $scope.showFull = function(self, next) {
    $('#hoverer img').attr({
      'src': self,
      'data-next': next
    });
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

  var lastScrollHeigth = document.body.scrollHeight;

  $(window).scroll(function() {

    var totalScroll = document.body.scrollHeight-$(window).height();
    var scrollTop = $(window).scrollTop();

    if(totalScroll/2 < scrollTop && lastScrollHeigth < document.body.scrollHeight) {
      $scope.nextPage();
      lastScrollHeigth = document.body.scrollHeight;
    }
  });

  $('#hoverer').on('click', function(e) {
    if(e.target.localName == "img") {
      $('#hover-img').attr('src', getNextImage($('#hover-img').attr('src'), $scope.images));
    } else {
      $('#hoverer').hide();
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

function getNextImage(currentSrc, all_images) {

  for(var key in all_images) {

    if(all_images[key].url == currentSrc) {
      return all_images[parseInt(key)+1].url;
    }
  }
  return currentSrc;
}