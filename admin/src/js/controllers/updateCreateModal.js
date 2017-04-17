angular.module('RDash')
	.controller('UpdateCreateModal', ['$scope', '$http', '$uibModalInstance', 'build', 'branch', UpdateCreateModal]);

function UpdateCreateModal($scope, $http, $uibModalInstance, build, branch) {
    $scope.build = build;
    $scope.branch = branch;
    $scope.update = {}
    $scope.update.log = build.Log;
    $scope.update.semver = "";
    $scope.close = function(){
        $uibModalInstance.close();
    }
    $scope.createUpdate = function(){
       $http.post(APIURL + '/addUpdate', {sessionid: $scope.Self.SessionId, semver: $scope.update.semver, build: build.Id, branch: branch.Name, log: $scope.update.log})
		.success(function () {
            if(!true){ // Needs to upload file too

                var fd = new FormData();
                fd.append('file', $scope.update.file);
                fd.append('sessionid', $scope.Self.SessionId);
                fd.append('build', build.Id);
                fd.append('branch', branch.Name);
                $http.post(APIURL + '/uploadUpdate', fd).success(function(){
                    $scope.refreshList();
			        $uibModalInstance.close();
                }).error(function(data, status){
                    $scope.parseError(data, status);
			        $scope.addAlert("Could not create update!", "danger");
                });
            } else {
                $scope.refreshList();
			    $uibModalInstance.close();
            }
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not create update!", "danger");
		});
    }
}