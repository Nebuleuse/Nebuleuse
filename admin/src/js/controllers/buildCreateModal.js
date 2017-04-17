angular.module('RDash')
	.controller('BuildCreateModal', ['$scope', '$http', '$uibModalInstance', 'list', 'commit', BuildCreateModal]);

function BuildCreateModal($scope, $http, $uibModalInstance, list, commit) {
	$scope.commits = [];
	$scope.diffs = [];
	$scope.list = list;
	$scope.fromCommit = commit;
	$scope.rawSize = 0;
	$scope.displaySize=0;
	$scope.showFiles = false;
	$scope.log = "";
	
	$http.post(APIURL + '/prepareGitBuild', {sessionid: $scope.Self.SessionId, commit: commit.Id})
		.success(function (data) {
			$scope.diffs = data.Diffs;
			$scope.rawSize = data.TotalSize;
			$scope.displaySize = Math.round(($scope.rawSize/1024) * 100) / 100;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch build preperation infos!", "danger");
		});

	var found = false;
	for (var i = 0; i < list.Commits.length; i++) {
		if (list.Commits[i].Id == commit.Id){
			found = true;
		}
		if (found){
			$scope.log += "-- Message from " + list.Commits[i].Id + "\n"+ list.Commits[i].Message; 
			$scope.commits.push(list.Commits[i]);
		}
	};
	
	$scope.toggleFiles = function () {
		$scope.showFiles = !$scope.showFiles;
	}
	$scope.close = function () {
		$uibModalInstance.close();
	}
	$scope.createBuild = function () {
		$http.post(APIURL + '/addGitBuild', {sessionid: $scope.Self.SessionId, commit: commit.Id, log: $scope.log})
		.success(function () {
			$scope.refreshList();
			$uibModalInstance.close();
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not create build!", "danger");
		});

	}
}