angular.module('RDash')
    .controller('LoginCtrl', ['$scope', '$cookieStore','$http', '$location', LoginCtrl]);

function LoginCtrl($scope, $cookieStore, $http, $location) {
	$scope.setPageTitle("Login");

	$scope.connect = function(username, password) {
        $http.post(APIURL + '/connect', {username: username, password: password})
        .success(function (data, status, headers, config) {
            $scope.Self.Username = username;
            $scope.Self.Password = password;
            $scope.Self.SessionId = data.SessionId;
            $cookieStore.put('sessionId', data.SessionId);
            $scope.getUserInfos();
            $scope.setConnected(true);
            $scope.getMessages();
            $location.path('/');
        }).error(function (data, status, headers, config) {
            $scope.parseError(data, status);
            $scope.addAlert("Impossible to login", 'error');
        });
    };
}