angular.module('RDash')
    .controller('ConfigCtrl', ['$scope', '$http', '$location', '$stateParams', ConfigCtrl]);

function ConfigCtrl($scope, $http, $location, $stateParams) {
	$scope.setPageTitle("Configuration");
	if(!$scope.checkAccess())
		return;
		
}
