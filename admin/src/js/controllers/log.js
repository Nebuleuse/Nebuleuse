angular.module('RDash')
    .controller('LogCtrl', ['$scope', '$http', LogCtrl]);

function LogCtrl($scope, $http) {
	if(!$scope.checkAccess())
		return;
	
	$scope.logLines = "";
	$scope.setPageTitle("Live Log");
	$scope.subscribeTo("admin", "log");

	$scope.$on("log", function (event, arg) {
		$scope.logLines += arg + "\n";
	});

	$http.post(APIURL + '/getLogs', {sessionid: $scope.Self.SessionId})
	.success(function (data) {
		$scope.logLines = data;
	})
	.error(function (data, status) {
		$scope.parseError(data, status);
		$scope.addAlert("Can't get past log");
    });


    $scope.clearLog = function() {
        $scope.logLines = "";
    };
}