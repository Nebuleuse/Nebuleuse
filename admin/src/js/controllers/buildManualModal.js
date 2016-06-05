angular.module('RDash')
	.controller('ManualBulidCreateModal', ['$scope', '$http', '$uibModalInstance', ManualBulidCreateModal]);

function ManualBulidCreateModal($scope, $http, $uibModalInstance, build, branch) {
	$scope.log = "";
    $scope.close = function(){
        $uibModalInstance.close();
    }
    $scope.createUpdate = function(){
        //todo
      /* $http.post(APIURL + '/uploadBuild', {sessionid: $scope.Self.SessionId})
		.success(function () {
			$scope.refreshList();
			$uibModalInstance.close();
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not create update!", "danger");
		});*/
    }
}